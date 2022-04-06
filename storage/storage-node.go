package storage

import (
	"context"
	"encoding/json"
	"math"
	"net"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/common/wait"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const SSTableRootDir = "data"

type StorageNode struct {
	*common.BaseComponent
	proto.UnimplementedKeyValueServer
	proto.UnimplementedStorageServer

	globalRwMutex     sync.RWMutex
	memTableRwMutexes map[common.SlotID]sync.RWMutex
	mem               *db.SlotMemTable[db.InternalKey, string]

	tableMgr  *TableManager // Init at Start()
	compactor *Compactor

	walMutex sync.Mutex
	wals     map[common.SlotID]*LogWriter
	store    *minio.Client

	// Init at Start()
	coord proto.CoordinatorClient
}

func NewStorageNode(addr string, etcd *etcdcli.Client, minio *minio.Client) *StorageNode {
	return &StorageNode{
		BaseComponent: common.NewBaseComponent(addr, etcd),

		globalRwMutex:     sync.RWMutex{},
		memTableRwMutexes: make(map[common.SlotID]sync.RWMutex),
		mem:               db.NewSlotMemTable[db.InternalKey, string](),

		compactor: NewCompactor(),

		walMutex: sync.Mutex{},
		wals:     make(map[common.SlotID]*LogWriter),
		store:    minio,
	}
}

func (node *StorageNode) Start() {
	log.Info("storage node starting...")

	ctx := context.Background()
	var coord proto.CoordinatorClient
	for i := 0; i < 3; i++ {
		resp, err := node.Meta.Get(ctx, "services/coordinator", etcdcli.WithPrefix())
		if err != nil {
			panic(err)
		}

		if resp.Count > 0 {
			var info proto.NodeInfo
			err = json.Unmarshal(resp.Kvs[0].Value, &info)
			if err != nil {
				panic(err)
			}

			conn, err := grpc.Dial(info.Addr, grpc.WithInsecure())
			if err != nil {
				panic(err)
			}

			coord = proto.NewCoordinatorClient(conn)
			log.Info("coordinator connected",
				zap.String("addr", info.Addr))
			break
		}

		time.Sleep(time.Second)
	}
	node.coord = coord
	node.tableMgr = NewTableManager("cyberkv", node.store, coord)

	listener, err := net.Listen("tcp", node.Info.Addr)
	if err != nil {
		panic(err)
	}

	go node.heartbeat()

	server := grpc.NewServer()
	proto.RegisterKeyValueServer(server, node)
	proto.RegisterStorageServer(server, node)
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func (node *StorageNode) Recovery() {
}

func (node *StorageNode) CompactMemTableAsLeader(ctx context.Context, request *proto.CompactMemTableRequest) (*proto.CompactMemTableResponse, error) {
	slot := common.SlotID(request.Slot)
	group := node.mem.Rotate(slot)
	imm := group.Imm()

	log.Info("ready to compact memtable as leader",
		zap.Int32("slot", slot))

	log.Info("create slot channel",
		zap.Int32("slot", slot))
	slotCh := node.compactor.PreCompact(slot)
	ch := slotCh.CreateDataChan()

	go func() {
		imm.Range(func(key db.InternalKey, value string) bool {
			ch <- &proto.KvData{
				Key:       key.UserKey(),
				Value:     value,
				Timestamp: key.GetTimeStamp(),
			}

			return true
		})

		close(ch)
	}()

	// Wait for the other storage nodes to push their memtables.
	isSatisfied := wait.WaitForCondition(10*time.Second, func() bool {
		return slotCh.Len() >= common.DefaultReadQuorum
	})

	if !isSatisfied {
		log.Error("failed to reach read quorum when compaction",
			zap.Int32("slot", slot),
			zap.Int("participantNumber", slotCh.Len()))
	}

	slotCh.Close()

	mergeCh := node.compactor.MergeChan(slotCh, math.MaxUint64)

	newSSTable := db.NewSSTableFromDataCh(mergeCh)
	err := node.tableMgr.WriteLevel0SSTable(ctx, newSSTable, 0)
	// created, deleted, err := node.tableMgr.WriteSSTable(ctx, mergeCh, 0)
	if err != nil {
		log.Error("failed to write sstable",
			zap.Int("level", 0),
			zap.Error(err))
		return nil, err
	}

	log.Info("compact done",
		zap.Int32("slot", slot))

	createdSSTables := make([]*proto.SSTableLevel, common.MaxLevel)
	deletedSSTables := make([]*proto.SSTableLevel, common.MaxLevel)

	for i := range createdSSTables {
		createdSSTables[i] = &proto.SSTableLevel{
			Level:    int32(i),
			Sstables: make([]string, 0),
		}
		deletedSSTables[i] = &proto.SSTableLevel{
			Level:    int32(i),
			Sstables: make([]string, 0),
		}
	}

	createdSSTables[0].Sstables = append(createdSSTables[0].Sstables, newSSTable.Path)
	// for _, sstable := range created {
	// 	createdSSTables[sstable.Level].Sstables = append(createdSSTables[sstable.Level].Sstables, sstable.Path)
	// }

	// for _, sstable := range deleted {
	// 	deletedSSTables[sstable.Level].Sstables = append(deletedSSTables[sstable.Level].Sstables, sstable.Path)
	// }

	return &proto.CompactMemTableResponse{
		CreatedSstables: createdSSTables,
		DeletedSstables: deletedSSTables,
	}, nil
}

func (node *StorageNode) CompactMemTableAsFollower(ctx context.Context, request *proto.CompactMemTableRequest) (*proto.CompactMemTableResponse, error) {
	slot := common.SlotID(request.Slot)
	group := node.mem.Rotate(slot)
	imm := group.Imm()

	log.Info("ready to compact memtable as follower",
		zap.Int32("slot", slot))

	data := GetKvDataBuffer()

	begin := time.Now()
	imm.Range(func(key db.InternalKey, value string) bool {
		data = append(data, &proto.KvData{
			Key:       key.UserKey(),
			Value:     value,
			Timestamp: key.GetTimeStamp(),
		})

		return true
	})
	end := time.Now()

	log.Info("read imm done",
		zap.Int32("slot", slot),
		zap.Int("len", len(data)),
		zap.Float64("seconds", end.Sub(begin).Seconds()))

	req := proto.PushMemTableRequest{
		Slot: request.Slot,
		Data: data,
	}

	conn, err := grpc.Dial(request.LeaderAddr, grpc.WithInsecure())
	if err != nil {
		log.Error("failed to connect to leader",
			zap.Int32("slot", slot),
			zap.String("leader_addr", request.LeaderAddr),
			zap.Error(err))
		return nil, err
	}
	cli := proto.NewStorageClient(conn)

	_, err = cli.PushMemTable(ctx, &req)
	PutBackToPool(data)
	if err != nil {
		log.Error("failed to push memtable to leader",
			zap.Int32("slot", slot),
			zap.String("leader_addr", request.LeaderAddr),
			zap.Error(err))
		return nil, err
	}

	return &proto.CompactMemTableResponse{}, nil
}

func (node *StorageNode) heartbeat() {
	ctx := context.Background()
	for {
		sizeTable := node.mem.GetSizeTable()
		if len(sizeTable) > 0 {
			node.coord.ReportStats(ctx, &proto.ReportStatsRequest{
				Id:           node.Info.Id,
				MemTableSize: sizeTable,
			})
		}

		time.Sleep(30 * time.Second)
	}
}

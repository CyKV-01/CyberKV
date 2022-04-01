package storage

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const SSTableRootDir = "data"

type StorageNode struct {
	*common.BaseNode
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
		BaseNode: common.NewBaseNode(addr, etcd),

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

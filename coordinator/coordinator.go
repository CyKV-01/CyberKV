package coordinator

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var _ common.Node = (*Coordinator)(nil)

type Coordinator struct {
	*common.BaseNode
	proto.UnimplementedKeyValueServer
	proto.UnimplementedCoordinatorServer

	computeCluster *Cluster[*ComputeNode]
	storageCluster *Cluster[*StorageNode]

	slotMemSizeTable []uint64

	compactor          *Compactor
	sstableIdAllocator *common.IdAllocator

	ts common.TimeStamp
}

func NewCoordinator(etcdClient *etcdcli.Client, addr string) *Coordinator {
	replicaNum := common.DefaultReplicaNum
	readQuorum := common.DefaultReadQuorum
	writeQuorum := common.DefaultWriteQuorum

	allocator, err := common.NewIdAllocator(context.Background(), etcdClient, common.SSTableIdKey, 100)
	if err != nil {
		panic(err)
	}

	return &Coordinator{
		BaseNode:       common.NewBaseNode(addr, etcdClient),
		computeCluster: NewCluster[*ComputeNode](etcdClient, 1, 1, 1),
		storageCluster: NewCluster[*StorageNode](etcdClient, replicaNum, readQuorum, writeQuorum),

		slotMemSizeTable:   make([]uint64, common.SlotNum),
		compactor:          NewCompactor(),
		sstableIdAllocator: allocator,
	}
}

func (coord *Coordinator) Start() {
	log.Info("coordinator starting...")

	coord.watchCluster()
	coord.Recovery()

	go coord.computeCluster.AssignSlotsBackground()
	go coord.storageCluster.AssignSlotsBackground()

	listener, err := net.Listen("tcp", coord.Info.Addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	proto.RegisterKeyValueServer(server, coord)
	proto.RegisterCoordinatorServer(server, coord)
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func (coord *Coordinator) Recovery() {
	// Recovery service info
	resp, err := coord.Meta.Get(context.Background(), common.ServicePrefix, etcdcli.WithPrefix())
	if err != nil {
		log.Errorf("failed to get cluster from etcd, err=%v", err)
		panic(err)
	}

	log.Infof("found %d nodes, recovering...", len(resp.Kvs))
	for _, kv := range resp.Kvs {
		coord.handleWatchEvent(kv)
	}

	coord.recoverySlotInfo()
}

func (coord *Coordinator) recoverySlotInfo() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)
	resp, err := coord.Meta.Get(ctx, common.SlotPrefix, etcdcli.WithPrefix())
	if err != nil {
		panic(err)
	}
	log.Infof("found %d slots, recovering...", len(resp.Kvs)/2)

	// etcd path: slots/{slot_id}/{compute/storage} -> NodeInfo
	computeSlots := make(map[common.SlotID]*proto.SlotInfo, common.SlotNum)
	storageSlots := make(map[common.SlotID]*proto.SlotInfo, common.SlotNum)
	for _, kv := range resp.Kvs {
		var info proto.SlotInfo

		slot := common.GetSlotIdByEtcdPath(string(kv.Key))

		err := json.Unmarshal(kv.Value, &info)
		if err != nil {
			panic(err)
		}

		if strings.Contains(string(kv.Key), "compute") {
			computeSlots[slot] = &info
		} else {
			storageSlots[slot] = &info
		}
	}

	for slot := int32(0); slot < common.SlotNum; slot++ {
		var (
			slotInfo *proto.SlotInfo
			ok       bool
		)

		if slotInfo, ok = computeSlots[common.SlotID(slot)]; !ok {
			slotInfo = &proto.SlotInfo{
				Slot:  slot,
				Nodes: make(map[string]*proto.NodeInfo),
			}
		}
		coord.computeCluster.RecoverySlotInfo(slotInfo)

		if slotInfo, ok = storageSlots[common.SlotID(slot)]; !ok {
			slotInfo = &proto.SlotInfo{
				Slot:  slot,
				Nodes: make(map[string]*proto.NodeInfo),
			}
		}
		coord.storageCluster.RecoverySlotInfo(slotInfo)
	}
}

func (coord *Coordinator) watchCluster() {
	log.Info("start watching cluster...")

	watchCh := coord.Meta.Watch(context.Background(), common.ServicePrefix, etcdcli.WithPrefix())

	go func() {
		for resp := range watchCh {
			for _, event := range resp.Events {
				if event.Type == etcdcli.EventTypePut {
					coord.handleWatchEvent(event.Kv)
				}
			}
		}
	}()
}

func (coord *Coordinator) handleWatchEvent(kv *mvccpb.KeyValue) {
	var nodeInfo VersionedNodeInfo
	nodeInfo.Id = string(kv.Key)
	nodeInfo.Addr = string(kv.Value)
	nodeInfo.Version = kv.CreateRevision

	if bytes.Contains(kv.Key, []byte("compute")) {
		log.Info("add new compute node",
			zap.String("id", nodeInfo.Id),
			zap.String("addr", nodeInfo.Addr))

		node, err := NewComputeNode(&nodeInfo)
		if err != nil {
			log.Errorf("failed to create node, err=v", err)
			return
		}
		coord.computeCluster.AddNode(node)
	} else if bytes.Contains(kv.Key, []byte("storage")) {
		err := json.Unmarshal(kv.Value, &nodeInfo)
		if err != nil {
			log.Error("failed to unmarshal storage node info",
				zap.Error(err))
			return
		}

		log.Info("add new storage node",
			zap.String("id", nodeInfo.Id),
			zap.String("addr", nodeInfo.Addr))

		node, err := NewStorageNode(&nodeInfo)
		if err != nil {
			log.Errorf("failed to create node, err=v", err)
			return
		}
		coord.storageCluster.AddNode(node)
	}
}

func (coord *Coordinator) GenTs() common.TimeStamp {
	return atomic.AddUint64(&coord.ts, 1)
}

func (coord *Coordinator) CurrentTs() common.TimeStamp {
	return atomic.LoadUint64(&coord.ts)
}

func (coord *Coordinator) AddSlotMemSize(slot common.SlotID, size uint64) uint64 {
	return atomic.AddUint64(&coord.slotMemSizeTable[slot], size)
}

func (coord *Coordinator) CompactMemTable(slot common.SlotID) {

}

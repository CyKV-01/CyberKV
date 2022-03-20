package coordinator

import (
	"bytes"
	"context"
	"encoding/json"
	"net"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Coordinator struct {
	*common.BaseNode
	proto.UnimplementedKeyValueServer
	proto.UnimplementedCoordinatorServer

	computeCluster *Cluster[ComputeNode]
	storageCluster *Cluster[StorageNode]
}

func NewCoordinator(etcdClient *etcdcli.Client, addr string) *Coordinator {
	replicaNum := DefaultReplicaNum
	readQuorum := DefaultReadQuorum
	writeQuorum := DefaultWriteQuorum

	return &Coordinator{
		BaseNode:       common.NewBaseNode(addr, etcdClient),
		computeCluster: NewCluster[ComputeNode](1, 1, 1),
		storageCluster: NewCluster[StorageNode](replicaNum, readQuorum, writeQuorum),
	}
}

func (coord *Coordinator) Start() {
	log.Info("coordinator starting...")

	go coord.watchCluster()
	go coord.computeCluster.AssignSlotsBackground()
	go coord.storageCluster.AssignSlotsBackground()

	listener, err := net.Listen("tcp", coord.Info.Addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	proto.RegisterKeyValueServer(server, coord)
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func (coord *Coordinator) watchCluster() {
	log.Info("start watching cluster...")

	watchCh := coord.Etcd.Watch(context.Background(), common.ServicePrefix, etcdcli.WithPrefix())

	resp, err := coord.Etcd.Get(context.Background(), common.ServicePrefix, etcdcli.WithPrefix())
	if err != nil {
		log.Errorf("failed to get cluster from etcd, err=%v", err)
		panic(err)
	}

	log.Infof("found %d nodes", len(resp.Kvs))
	for _, kv := range resp.Kvs {
		coord.handleWatchEvent(kv)
	}

	for resp := range watchCh {
		for _, event := range resp.Events {
			if event.Type == etcdcli.EventTypePut {
				coord.handleWatchEvent(event.Kv)
			}
		}
	}
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

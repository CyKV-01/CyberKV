package coordinator

import (
	"context"
	"encoding/json"
	"net"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Coordinator struct {
	proto.UnimplementedKeyValueServer
	proto.UnimplementedCoordinatorServer

	id   uuid.UUID
	info *proto.NodeInfo

	etcdClient *etcdcli.Client

	computeCluster Cluster
	storageCluster Cluster
}

func NewCoordinator(etcdClient *etcdcli.Client, addr string) *Coordinator {
	replicaNum := DefaultReplicaNum
	readQuorum := DefaultReadQuorum
	writeQuorum := DefaultWriteQuorum

	nodeInfo := proto.NodeInfo{
		Addr: addr,
	}

	return &Coordinator{
		id:             uuid.New(),
		info:           &nodeInfo,
		etcdClient:     etcdClient,
		computeCluster: NewComputeCluster(1, 1, 1),
		storageCluster: NewStorageCluster(replicaNum, readQuorum, writeQuorum),
	}
}

func (coord *Coordinator) Start() {
	coord.Register()
	go coord.watchCluster()

	listener, err := net.Listen("tcp", coord.info.Addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	proto.RegisterKeyValueServer(server, coord)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

func (coord *Coordinator) Register() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)

	resp, err := coord.etcdClient.Grant(ctx, common.DefaultTTL)
	if err != nil {
		panic(err)
	}

	keepAliveCh, err := coord.etcdClient.KeepAlive(ctx, resp.ID)
	if err != nil {
		panic(err)
	}

	go func() {
		select {
		case <-keepAliveCh:
			return
		}
	}()

	infoBytes, err := json.Marshal(coord.info)
	if err != nil {
		panic(err)
	}

	log.Info("register coordinator",
		zap.String("id", coord.id.String()),
		zap.String("addr", coord.info.Addr))
	_, err = coord.etcdClient.Put(ctx, path.Join(common.ServicePrefix, "coordinator", coord.id.String()), string(infoBytes))
	if err != nil {
		panic(err)
	}
}

func (coord *Coordinator) watchCluster() {
	watchCh := coord.etcdClient.Watch(context.Background(), common.ServicePrefix)

	for resp := range watchCh {
		for _, event := range resp.Events {
			if event.Type == etcdcli.EventTypePut {
				var nodeInfo proto.NodeInfo
				err := json.Unmarshal(event.Kv.Value, &nodeInfo)
				if err != nil {
					log.Info("failed to watch new node regisering", zap.Error(err))
				} else {
					if nodeInfo.Type == proto.NodeType_ComputeNode {
						coord.computeCluster.AddNode(&nodeInfo)
					} else if nodeInfo.Type == proto.NodeType_StorageNode {
						coord.storageCluster.AddNode(&nodeInfo)
					}
				}
			}
		}
	}
}

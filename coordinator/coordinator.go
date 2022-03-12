package coordinator

import (
	"context"
	"encoding/json"
	"log"
	"path"

	"github.com/google/uuid"
	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
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
	writeQuorum := DefaultWriteQuorum
	readQuorum := DefaultReadQuorum

	nodeInfo := proto.NodeInfo{
		Addr: addr,
	}

	return &Coordinator{
		id:             uuid.New(),
		info:           &nodeInfo,
		etcdClient:     etcdClient,
		computeCluster: NewComputeCluster(1, 1, 1),
		storageCluster: NewStorageCluster(replicaNum, writeQuorum, readQuorum),
	}
}

func (coord *Coordinator) Start() {
	coord.Register()

	go coord.watchCluster()
}

func (coord *Coordinator) Register() {
	ctx := context.Background()

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
					log.Printf("failed to watch new node regisering, err=%+v", err)
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

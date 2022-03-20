package storage

import (
	"net"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const SSTableRootDir = "data"

type StorageNode struct {
	*common.BaseNode
	proto.UnimplementedKeyValueServer

	walMutex sync.Mutex
	wals     map[common.SlotID]*LogWriter
	store    *minio.Client
}

func NewStorageNode(addr string, etcd *etcdcli.Client, minio *minio.Client) *StorageNode {
	return &StorageNode{
		BaseNode: common.NewBaseNode(addr, etcd),
		walMutex: sync.Mutex{},
		wals:     make(map[common.SlotID]*LogWriter),
		store:    minio,
	}
}

func (node *StorageNode) Start() {
	log.Info("coordinator starting...")

	listener, err := net.Listen("tcp", node.Info.Addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	proto.RegisterKeyValueServer(server, node)
	reflection.Register(server)
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

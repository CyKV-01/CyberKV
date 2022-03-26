package storage

import (
	"net"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
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

	globalRwMutex     sync.RWMutex
	memTableRwMutexes map[common.SlotID]sync.RWMutex
	mem               *db.SlotMemTable[db.InternalKey, string]
	// imm               *db.SlotMemTable[db.InternalKey, string]
	// fmem              *db.SlotMemTable[db.InternalKey, string]

	walMutex sync.Mutex
	wals     map[common.SlotID]*LogWriter
	store    *minio.Client
}

func NewStorageNode(addr string, etcd *etcdcli.Client, minio *minio.Client) *StorageNode {
	return &StorageNode{
		BaseNode: common.NewBaseNode(addr, etcd),

		globalRwMutex:     sync.RWMutex{},
		memTableRwMutexes: make(map[int16]sync.RWMutex),
		mem:               db.NewSlotMemTable[db.InternalKey, string](),
		// imm:               nil,
		// fmem:              nil,

		walMutex: sync.Mutex{},
		wals:     make(map[common.SlotID]*LogWriter),
		store:    minio,
	}
}

func (node *StorageNode) Start() {
	log.Info("storage node starting...")

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

func (node *StorageNode) CompactMemTable() {
	// node.rotateMemTables()
}
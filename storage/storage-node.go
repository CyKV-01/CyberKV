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

	memTableSwitchRWMutex sync.RWMutex
	mem                   *db.SlotMemTable[db.InternalKey, string]
	imm                   *db.SlotMemTable[db.InternalKey, string]
	fmem                  *db.SlotMemTable[db.InternalKey, string]

	walMutex sync.Mutex
	wals     map[common.SlotID]*LogWriter
	store    *minio.Client
}

func NewStorageNode(addr string, etcd *etcdcli.Client, minio *minio.Client) *StorageNode {
	return &StorageNode{
		BaseNode: common.NewBaseNode(addr, etcd),

		memTableSwitchRWMutex: sync.RWMutex{},
		mem:                   db.NewSlotMemTable[db.InternalKey, string](),
		imm:                   nil,
		fmem:                  nil,

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

}

func (node *StorageNode) getMemTables() (*db.SlotMemTable[db.InternalKey, string], *db.SlotMemTable[db.InternalKey, string], *db.SlotMemTable[db.InternalKey, string]) {
	node.memTableSwitchRWMutex.RLock()
	defer node.memTableSwitchRWMutex.RUnlock()

	return node.mem, node.imm, node.fmem
}

func (node *StorageNode) rotateMemTables() {
	node.memTableSwitchRWMutex.Lock()
	defer node.memTableSwitchRWMutex.Unlock()

	node.fmem = node.imm
	node.imm = node.mem
	node.mem = db.NewSlotMemTable[db.InternalKey, string]()
}

package coordinator

import (
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
	"google.golang.org/grpc"
)

type VersionedNodeInfo struct {
	proto.NodeInfo
	Version int64
}

type Node interface {
	// GetId() common.NodeID
	GetSlots() []common.SlotID
	// HasSlot(id common.SlotID) bool
	// GetUsage() uint64
	// GetPercent() float32
	// GetCap() uint64
	// GetCpuPercent() float32

	AssignSlots(slots []common.SlotID) error

	GetInfo() *VersionedNodeInfo
}

type baseNode struct {
	info *VersionedNodeInfo

	rwmutex             sync.RWMutex // guard fields below
	serveSlots          map[common.SlotID]struct{}
	storageUsage        uint64
	storageUsagePercent float32
	cap                 uint64 // not guard, cap field is read-only
	cpuPercent          float32

	proto.KeyValueClient
}

func NewBaseNode(info *VersionedNodeInfo) (*baseNode, error) {
	conn, err := grpc.Dial(info.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	kvClient := proto.NewKeyValueClient(conn)

	return &baseNode{
		info:           info,
		rwmutex:        sync.RWMutex{},
		serveSlots:     make(map[common.SlotID]struct{}),
		KeyValueClient: kvClient,
	}, nil
}

func (node *baseNode) GetSlots() []common.SlotID {
	node.rwmutex.RLock()
	defer node.rwmutex.RUnlock()

	slots := make([]common.SlotID, 0, len(node.serveSlots))
	for id := range node.serveSlots {
		slots = append(slots, id)
	}

	return slots
}

func (node *baseNode) AssignSlots(slots []common.SlotID) error {
	node.rwmutex.Lock()
	defer node.rwmutex.Unlock()

	for _, slot := range slots {
		node.serveSlots[slot] = struct{}{}
	}

	return nil
}

func (node *baseNode) GetInfo() *VersionedNodeInfo {
	return node.info
}

type ComputeNode struct {
	*baseNode
}

func NewComputeNode(info *VersionedNodeInfo) (ComputeNode, error) {
	baseNode, err := NewBaseNode(info)
	if err != nil {
		return ComputeNode{}, err
	}
	node := ComputeNode{
		baseNode: baseNode,
	}
	return node, nil
}

type StorageNode struct {
	*baseNode
}

func NewStorageNode(info *VersionedNodeInfo) (StorageNode, error) {
	baseNode, err := NewBaseNode(info)
	if err != nil {
		return StorageNode{}, err
	}
	node := StorageNode{
		baseNode: baseNode,
	}
	return node, nil
}

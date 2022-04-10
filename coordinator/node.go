package coordinator

import (
	"context"
	"sync"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type VersionedNodeInfo struct {
	proto.NodeInfo
	Version int64
}

type NodeType int

const (
	ComputeNodeType NodeType = iota + 1
	StorageNodeType
)

type Node interface {
	// GetId() common.NodeID
	GetSlots() []common.SlotID
	// HasSlot(id common.SlotID) bool
	// GetUsage() uint64
	// GetPercent() float32
	// GetCap() uint64
	// GetCpuPercent() float32
	GetNodeType() NodeType

	HasSlot(common.SlotID) bool

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
		log.Error("failed to dial node",
			zap.Uint64("nodeID", info.Id),
			zap.String("addr", info.Addr),
			zap.Error(err))
		return nil, err
	}

	kvClient := proto.NewKeyValueClient(conn)

	// _, err = kvClient.AssignSlot(context.Background(), &proto.AssignSlotRequest{SlotID: 0})
	// if err != nil {
	// 	panic(err)
	// }

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
		ctx := context.Background()
		time.Sleep(time.Second)
		_, err := node.AssignSlot(ctx, &proto.AssignSlotRequest{SlotID: slot})
		if err != nil {
			return err
		}
	}

	return nil
}

func (node *baseNode) GetInfo() *VersionedNodeInfo {
	return node.info
}

func (node *baseNode) HasSlot(slot common.SlotID) bool {
	node.rwmutex.RLock()
	defer node.rwmutex.RUnlock()

	_, ok := node.serveSlots[slot]

	return ok
}

type ComputeNode struct {
	*baseNode
}

func NewComputeNode(info *VersionedNodeInfo) (*ComputeNode, error) {
	baseNode, err := NewBaseNode(info)
	if err != nil {
		return nil, err
	}
	node := ComputeNode{
		baseNode: baseNode,
	}
	return &node, nil
}

func (node *ComputeNode) GetNodeType() NodeType {
	return ComputeNodeType
}

type StorageNode struct {
	*baseNode
	proto.StorageClient
}

func NewStorageNode(info *VersionedNodeInfo) (*StorageNode, error) {
	baseNode, err := NewBaseNode(info)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(info.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	storageClient := proto.NewStorageClient(conn)

	node := StorageNode{
		baseNode:      baseNode,
		StorageClient: storageClient,
	}
	return &node, nil
}

func (node *StorageNode) GetNodeType() NodeType {
	return StorageNodeType
}

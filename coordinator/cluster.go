package coordinator

import (
	"sync"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
)

const (
	MaxNodeUsagePercent = 0.8
	DefaultReplicaNum   = 3
	DefaultWriteQuorum  = 2
	DefaultReadQuorum   = 2
)

type SlotInfo struct {
	nodes []Node
}

type Cluster interface {
	AddNode(*proto.NodeInfo)
	// AssignSlots(slots []common.SlotID, excludedNodes []common.NodeID) (int, error)
	// GetNodes(excludedNodes []common.NodeID) []Node
	// GetNode(id common.NodeID) Node
	GetNodesBySlot(slot common.SlotID) []Node
	// AssignSlotsBackground()
}

type Node interface {
	// GetId() common.NodeID
	// GetSlots() []common.SlotID
	// HasSlot(id common.SlotID) bool
	// GetUsage() uint64
	// GetPercent() float32
	// GetCap() uint64
	// GetCpuPercent() float32

	// AssignSlots(slots []common.SlotID) error
}

type baseCluster struct {
	slots       map[common.SlotID]*SlotInfo
	replicaNum  int
	readQuorum  int
	writeQuorum int

	nodes   map[common.NodeID]Node
	rwmutex sync.RWMutex
}

func (cluster *baseCluster) AddNode(info *proto.NodeInfo) {
	cluster.rwmutex.Lock()
	defer cluster.rwmutex.Unlock()

	cluster.nodes[common.NodeID(info.Id)] = &baseNode{
		info: info,
	}
}

// func (cluster *baseCluster) AssignSlots(slots []common.SlotID, excludedNodes []common.NodeID) (int, error) {
// 	nodes := cluster.GetNodes(excludedNodes)
// 	sort.Slice(nodes, func(i, j int) bool {
// 		return len(nodes[i].GetSlots()) < len(nodes[j].GetSlots())
// 	})

// 	if len(nodes) == 0 {
// 		return 0, fmt.Errorf("no node to assign slots")
// 	}

// 	nodeIdx := 0
// 	j := 0
// 	for i := 0; i < len(slots); i = j {
// 		node := nodes[nodeIdx]
// 		j := i + 3
// 		if j > len(slots) {
// 			j = len(slots)
// 		}

// 		err := node.AssignSlots(slots[i:j])
// 		if err != nil {
// 			return i, err
// 		}
// 	}

// 	return len(slots), nil
// }

// func (cluster *baseCluster) GetNodes(excludedNodes []common.NodeID) []Node {
// 	excludedNodesMap := make(map[common.NodeID]bool, len(excludedNodes))
// 	for _, node := range excludedNodes {
// 		excludedNodesMap[node] = true
// 	}

// 	cluster.rwmutex.RLock()
// 	defer cluster.rwmutex.RUnlock()

// 	nodes := make([]Node, 0, len(cluster.nodes)-len(excludedNodes))
// 	for id, node := range cluster.nodes {
// 		if !excludedNodesMap[id] {
// 			nodes = append(nodes, node)
// 		}
// 	}

// 	return nodes
// }

// func (cluster *baseCluster) GetNode(id common.NodeID) Node {
// 	cluster.rwmutex.RLock()
// 	defer cluster.rwmutex.RUnlock()

// 	return cluster.nodes[id]
// }

func (cluster *baseCluster) GetNodesBySlot(slot common.SlotID) []Node {
	slotInfo, ok := cluster.slots[slot]
	if !ok {
		return nil
	}

	return slotInfo.nodes
}

func (cluster *baseCluster) AssignSlotsBackground() {
	for ; true; time.Sleep(time.Second) {

		slots := cluster.getUnassignedSlots()
		excludedNodes := []common.NodeID{}
		slotIDs := []common.SlotID{}

		for _, slot := range slots {
			cluster.AssignSlots()
		}
	}
}

func (cluster *baseCluster) getUnassignedSlots() []*SlotInfo {
	unassignedSlots := make([]*SlotInfo, 0)

	for _, slot := range cluster.slots {
		if len(slot.nodes) < cluster.replicaNum {
			unassignedSlots = append(unassignedSlots, slot)
		}
	}

	return unassignedSlots
}

type baseNode struct {
	info *proto.NodeInfo

	serveSlots          map[common.SlotID]*SlotInfo
	storageUsage        uint64
	storageUsagePercent float32
	cap                 uint64
	cpuPercent          float32
}

type ComputeCluster struct {
	*baseCluster
}

func NewComputeCluster(replicaNum, writeQuorum, readQuorum int) *ComputeCluster {
	return &ComputeCluster{
		baseCluster: &baseCluster{
			slots: make(map[common.SlotID]*SlotInfo),

			replicaNum:  replicaNum,
			writeQuorum: writeQuorum,
			readQuorum:  readQuorum,

			nodes:   make(map[common.NodeID]Node),
			rwmutex: sync.RWMutex{},
		},
	}
}

type ComputeNode struct {
	*baseNode
	proto.KeyValueClient
}

type StorageCluster struct {
	*baseCluster
}

func NewStorageCluster(replicaNum, writeQuorum, readQuorum int) *StorageCluster {
	return &StorageCluster{
		baseCluster: &baseCluster{
			slots: make(map[common.SlotID]*SlotInfo),

			replicaNum:  replicaNum,
			writeQuorum: writeQuorum,
			readQuorum:  readQuorum,

			nodes:   make(map[common.NodeID]Node),
			rwmutex: sync.RWMutex{},
		},
	}
}

type StorageNode struct {
	*baseNode
}

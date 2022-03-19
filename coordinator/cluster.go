package coordinator

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
)

const (
	MaxNodeUsagePercent = 0.8
	DefaultReplicaNum   = 3
	DefaultWriteQuorum  = 2
	DefaultReadQuorum   = 2
)

type SlotInfo struct {
	Id    common.SlotID
	Nodes []Node
}

type Cluster interface {
	AddNode(*VersionedNodeInfo)
	AssignSlotsBackground()
	AssignSlots(slots []*SlotInfo) (int, error)
	GetNodes() []Node
	// GetNode(id common.NodeID) Node
	GetNodesBySlot(slot common.SlotID) []Node
	// AssignSlotsBackground()
}

type baseCluster struct {
	slots            map[common.SlotID]*SlotInfo
	scheduleTimer    time.Timer
	scheduleNotifier chan []*SlotInfo

	replicaNum  int
	readQuorum  int
	writeQuorum int

	nodes   map[common.NodeID]Node
	rwmutex sync.RWMutex
}

func NewBaseCluster(replicaNum, readQuorum, writeQuorum int) *baseCluster {
	return &baseCluster{
		slots:            make(map[common.SlotID]*SlotInfo),
		scheduleTimer:    *time.NewTimer(time.Second),
		scheduleNotifier: make(chan []*SlotInfo, 32),

		replicaNum:  replicaNum,
		readQuorum:  readQuorum,
		writeQuorum: writeQuorum,

		nodes:   make(map[common.NodeID]Node),
		rwmutex: sync.RWMutex{},
	}
}

func (cluster *baseCluster) AddNode(info *VersionedNodeInfo) {
	cluster.rwmutex.Lock()
	defer cluster.rwmutex.Unlock()

	old, ok := cluster.nodes[common.NodeID(info.Id)]
	if !ok || old.GetInfo().Version < info.Version {
		cluster.nodes[common.NodeID(info.Id)] = &baseNode{
			info: info,
		}
	}
}

func (cluster *baseCluster) AssignSlots(slots []*SlotInfo) (int, error) {
	nodes := cluster.GetNodes()
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].GetSlots()) < len(nodes[j].GetSlots())
	})

	if len(nodes) == 0 {
		return 0, fmt.Errorf("no node to assign slots")
	}

	nodeIdx := 0
	j := 0
	assignedSlotCount := 0
	for i := 0; i < len(slots) && i < len(nodes); i = j {
		node := nodes[nodeIdx]
		j := i + 1
		if j > len(slots) {
			j = len(slots)
		}

		err := node.AssignSlots(slots[i:j])
		if err != nil {
			return i, err
		}
		assignedSlotCount += j - i
	}

	return assignedSlotCount, nil
}

func (cluster *baseCluster) GetNodes() []Node {
	// excludedNodesMap := make(map[common.NodeID]bool, len(excludedNodes))
	// for _, node := range excludedNodes {
	// 	excludedNodesMap[node] = true
	// }

	cluster.rwmutex.RLock()
	defer cluster.rwmutex.RUnlock()

	nodes := make([]Node, 0, len(cluster.nodes))
	for _, node := range cluster.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

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

	return slotInfo.Nodes
}

func (cluster *baseCluster) AssignSlotsBackground() {
	var slots []*SlotInfo
	for ; true; time.Sleep(time.Second) {
		select {
		case <-cluster.scheduleTimer.C:
			slots = cluster.getUnassignedSlots()

		case slots = <-cluster.scheduleNotifier:
		}

		if len(slots) == 0 {
			continue
		}

		log.Infof("schedule slots=%+v", slots)
		n, err := cluster.AssignSlots(slots)
		if n > 0 {
			log.Info("scheduled slots", zap.Int("succeed", n))
		}
		if err != nil {
			log.Error("failed to schedule some slots", zap.Error(err))
		}
	}
}

// Only created slots
func (cluster *baseCluster) getUnassignedSlots() []*SlotInfo {
	unassignedSlots := make([]*SlotInfo, 0)

	for _, slot := range cluster.slots {
		if len(slot.Nodes) < cluster.replicaNum {
			unassignedSlots = append(unassignedSlots, slot)
		}
	}

	return unassignedSlots
}

type ComputeCluster struct {
	*baseCluster
}

func NewComputeCluster(replicaNum, readQuorum, writeQuorum int) *ComputeCluster {
	return &ComputeCluster{
		baseCluster: NewBaseCluster(replicaNum, readQuorum, writeQuorum),
	}
}

type ComputeNode struct {
	*baseNode
	proto.KeyValueClient
}

type StorageCluster struct {
	*baseCluster
}

func NewStorageCluster(replicaNum, readQuorum, writeQuorum int) *StorageCluster {
	return &StorageCluster{
		baseCluster: NewBaseCluster(replicaNum, readQuorum, writeQuorum),
	}
}

type StorageNode struct {
	*baseNode
}

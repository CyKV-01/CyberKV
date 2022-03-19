package coordinator

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"go.uber.org/zap"
)

const (
	MaxNodeUsagePercent = 0.8
	DefaultReplicaNum   = 3
	DefaultWriteQuorum  = 2
	DefaultReadQuorum   = 2
)

type SlotInfo[T Node] struct {
	Id    common.SlotID
	Nodes []T
}

type ClusterI[T Node] interface {
	AddNode(*VersionedNodeInfo)
	AssignSlotsBackground()
	AssignSlots(slots []*SlotInfo[T])
	GetNodes() []T
	// GetNode(id common.NodeID) Node
	GetNodesBySlot(slot common.SlotID) []T // return nil if no such slot
	// AssignSlotsBackground()
}

type Cluster[T Node] struct {
	slots            map[common.SlotID]*SlotInfo[T]
	scheduleTimer    time.Timer
	scheduleNotifier chan []common.SlotID

	replicaNum  int
	readQuorum  int
	writeQuorum int

	nodes   map[common.NodeID]T
	rwmutex sync.RWMutex
}

func NewCluster[T Node](replicaNum, readQuorum, writeQuorum int) *Cluster[T] {
	return &Cluster[T]{
		slots:            make(map[common.SlotID]*SlotInfo[T]),
		scheduleTimer:    *time.NewTimer(time.Second),
		scheduleNotifier: make(chan []common.SlotID, 32),

		replicaNum:  replicaNum,
		readQuorum:  readQuorum,
		writeQuorum: writeQuorum,

		nodes:   make(map[common.NodeID]T),
		rwmutex: sync.RWMutex{},
	}
}

func (cluster *Cluster[T]) AddNode(node T) {
	cluster.rwmutex.Lock()
	defer cluster.rwmutex.Unlock()

	info := node.GetInfo()
	old, ok := cluster.nodes[common.NodeID(info.Id)]
	if !ok || old.GetInfo().Version < info.Version {
		cluster.nodes[common.NodeID(info.Id)] = node
	}
}

func (cluster *Cluster[T]) AssignSlots(slots []common.SlotID) {
	cluster.scheduleNotifier <- slots
}

func (cluster *Cluster[T]) assignSlots(slots []common.SlotID) (int, error) {
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
		j = i + 1
		if j > len(slots) {
			j = len(slots)
		}

		log.Infof("schedule slots=%+v to node=%+v", slots, node)
		scheduleSlots := slots[i:j]
		err := node.AssignSlots(scheduleSlots)
		if err != nil {
			return i, err
		}
		for _, slot := range scheduleSlots {
			slotInfo, ok := cluster.slots[slot]
			if !ok {
				slotInfo = &SlotInfo[T]{
					Id:    slot,
					Nodes: make([]T, 0, 1),
				}
				cluster.slots[slot] = slotInfo
			}

			slotInfo.Nodes = append(slotInfo.Nodes, node)
		}

		assignedSlotCount += j - i
	}

	return assignedSlotCount, nil
}

func (cluster *Cluster[T]) GetNodes() []T {
	// excludedNodesMap := make(map[common.NodeID]bool, len(excludedNodes))
	// for _, node := range excludedNodes {
	// 	excludedNodesMap[node] = true
	// }

	cluster.rwmutex.RLock()
	defer cluster.rwmutex.RUnlock()

	nodes := make([]T, 0, len(cluster.nodes))
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

func (cluster *Cluster[T]) GetNodesBySlot(slot common.SlotID) []T {
	slotInfo, ok := cluster.slots[slot]
	if !ok {
		return nil
	}

	log.Infof("slot=%+v", slotInfo)

	return slotInfo.Nodes
}

func (cluster *Cluster[T]) AssignSlotsBackground() {
	var slots []common.SlotID
	for {
		select {
		case <-cluster.scheduleTimer.C:
			slots = cluster.getUnassignedSlots()
			cluster.scheduleTimer.Reset(time.Second)

		case slots = <-cluster.scheduleNotifier:
		}

		if len(slots) == 0 {
			continue
		}

		log.Infof("schedule slots=%+v", slots)
		n, err := cluster.assignSlots(slots)
		if n > 0 {
			log.Info("scheduled slots", zap.Int("succeed", n))
		}
		if err != nil {
			log.Error("failed to schedule some slots", zap.Error(err))
		}
	}
}

// Only created slots
func (cluster *Cluster[T]) getUnassignedSlots() []common.SlotID {
	unassignedSlots := make([]common.SlotID, 0)

	for slot, info := range cluster.slots {
		if len(info.Nodes) < cluster.replicaNum {
			unassignedSlots = append(unassignedSlots, slot)
		}
	}

	return unassignedSlots
}

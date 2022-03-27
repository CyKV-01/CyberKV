package coordinator

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const (
	MaxNodeUsagePercent = 0.8
)

type SlotInfo[T Node] struct {
	Id    common.SlotID
	Nodes []T
}

// type ClusterI[T Node] interface {
// 	AddNode(*VersionedNodeInfo)
// 	AssignSlotsBackground()
// 	AssignSlots(slots []*SlotInfo[T])
// 	GetNodes() []T
// 	// GetNode(id common.NodeID) Node
// 	GetNodesBySlot(slot common.SlotID) []T // return nil if no such slot
// 	// AssignSlotsBackground()
// }

type Cluster[T Node] struct {
	meta        *etcdcli.Client
	replicaNum  int
	readQuorum  int
	writeQuorum int

	rwmutex sync.RWMutex
	slots   []*proto.SlotInfo
	nodes   map[common.NodeID]T

	scheduleTimer *time.Ticker
}

func NewCluster[T Node](meta *etcdcli.Client, replicaNum, readQuorum, writeQuorum int) *Cluster[T] {
	cluster := &Cluster[T]{
		meta:        meta,
		replicaNum:  replicaNum,
		readQuorum:  readQuorum,
		writeQuorum: writeQuorum,

		rwmutex: sync.RWMutex{},
		slots:   make([]*proto.SlotInfo, common.SlotNum),
		nodes:   make(map[common.NodeID]T),

		scheduleTimer: time.NewTicker(time.Second),
	}

	for i := 0; i < common.SlotNum; i++ {
		slot := common.SlotID(i)
		cluster.slots[slot] = &proto.SlotInfo{
			Slot:  uint32(slot),
			Nodes: make(map[string]*proto.NodeInfo, replicaNum),
		}
	}

	return cluster
}

func (cluster *Cluster[T]) AddNode(node T) {
	cluster.rwmutex.Lock()
	defer cluster.rwmutex.Unlock()

	info := node.GetInfo()
	old, ok := cluster.nodes[info.Id]
	if !ok || old.GetInfo().Version < info.Version {
		cluster.nodes[info.Id] = node
	}
}

func (cluster *Cluster[T]) RecoverySlotInfo(slotInfo *proto.SlotInfo) {
	recoveredNodes := make([]*proto.NodeInfo, 0, len(slotInfo.Nodes))
	for id, info := range slotInfo.Nodes {
		node, ok := cluster.GetNode(id)
		if !ok {
			continue
		}

		err := node.AssignSlots([]common.SlotID{common.SlotID(slotInfo.Slot)})
		if err != nil {
			log.Warn("failed to assign slot to node",
				zap.Error(err))
			continue
		}

		recoveredNodes = append(recoveredNodes, info)
	}

	slotInfo.Nodes = make(map[string]*proto.NodeInfo, len(recoveredNodes))
	for _, node := range recoveredNodes {
		slotInfo.Nodes[node.Id] = node
	}

	cluster.slots[slotInfo.Slot] = slotInfo
}

func (cluster *Cluster[T]) GetSlotInfo(id common.SlotID) *proto.SlotInfo {
	cluster.rwmutex.RLock()
	defer cluster.rwmutex.RUnlock()

	return cluster.slots[id]
}

func (cluster *Cluster[T]) assignSlots(slots []common.SlotID) error {
	nodes := cluster.GetNodes()
	if len(nodes) == 0 {
		return fmt.Errorf("no node to assign slots")
	}

	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].GetSlots()) < len(nodes[j].GetSlots())
	})

	nodeIdx := 0
	newSlotInfos := make(map[common.SlotID]*proto.SlotInfo, len(slots))
	for _, slot := range slots {
		newInfo, ok := newSlotInfos[slot]
		// Clone the old SlotInfo
		if !ok {
			old := cluster.GetSlotInfo(slot)
			newInfo = &proto.SlotInfo{
				Slot:  uint32(slot),
				Nodes: make(map[string]*proto.NodeInfo, len(old.Nodes)),
			}
			for id, info := range old.Nodes {
				newInfo.Nodes[id] = info
			}
			newSlotInfos[slot] = newInfo
		}

		for len(newInfo.Nodes) < cluster.replicaNum {
			node := nodes[nodeIdx]
			nodeIdx++
			if nodeIdx >= len(nodes) {
				nodeIdx %= len(nodes)
			}
			newInfo.Nodes[node.GetInfo().Id] = &node.GetInfo().NodeInfo
		}
	}

	nodeType := nodes[0].GetNodeType()
	wg := sync.WaitGroup{}
	for _, info := range newSlotInfos {
		wg.Add(1)
		go func(info *proto.SlotInfo) {
			defer wg.Done()

			err := cluster.saveSlotInfo(info, nodeType)
			if err != nil {
				return
			}

			succeedNodes := make([]*proto.NodeInfo, 0, len(info.Nodes))
			for id, nodeInfo := range info.Nodes {
				node, ok := cluster.GetNode(id)
				if !ok {
					continue
				}

				if !node.HasSlot(common.SlotID(info.Slot)) {
					err = node.AssignSlots([]common.SlotID{common.SlotID(info.Slot)})
					if err != nil {
						continue
					}
				}

				succeedNodes = append(succeedNodes, nodeInfo)
			}

			if len(succeedNodes) != len(info.Nodes) {
				info.Nodes = make(map[string]*proto.NodeInfo, len(succeedNodes))
				for _, nodeInfo := range succeedNodes {
					info.Nodes[nodeInfo.Id] = nodeInfo
				}

				err = cluster.saveSlotInfo(info, nodeType)
				if err != nil {
					return
				}
			}

		}(info)
	}

	wg.Wait()

	for id, info := range newSlotInfos {
		cluster.slots[id] = info
	}

	return nil
}

func (cluster *Cluster[T]) saveSlotInfo(slot *proto.SlotInfo, nodeType NodeType) error {
	nodeTypeStr := "compute"
	if nodeType == StorageNodeType {
		nodeTypeStr = "storage"
	}

	key := fmt.Sprintf("slots/%v/%v", slot.Slot, nodeTypeStr)
	value, err := json.Marshal(slot)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second)
	_, err = cluster.meta.Put(ctx, key, string(value))
	if err != nil {
		return err
	}

	return nil
}

func (cluster *Cluster[T]) GetNodes() []T {
	cluster.rwmutex.RLock()
	defer cluster.rwmutex.RUnlock()

	nodes := make([]T, 0, len(cluster.nodes))
	for _, node := range cluster.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

func (cluster *Cluster[T]) GetNode(id common.NodeID) (T, bool) {
	cluster.rwmutex.RLock()
	defer cluster.rwmutex.RUnlock()

	node, ok := cluster.nodes[id]

	return node, ok
}

// func (cluster *baseCluster) GetNode(id common.NodeID) Node {
// 	cluster.rwmutex.RLock()
// 	defer cluster.rwmutex.RUnlock()

// 	return cluster.nodes[id]
// }

func (cluster *Cluster[T]) GetNodesBySlot(slot common.SlotID) []T {
	slotInfo := cluster.slots[slot]

	log.Infof("slot=%+v", slotInfo)
	nodes := make([]T, 0, len(slotInfo.Nodes))
	for _, nodeInfo := range slotInfo.Nodes {
		nodes = append(nodes, cluster.nodes[nodeInfo.Id])
	}

	return nodes
}

func (cluster *Cluster[T]) AssignSlotsBackground() {
	var slots []common.SlotID
	for ; true; <-cluster.scheduleTimer.C {
		slots = cluster.getUnassignedSlots()
		if len(slots) == 0 {
			continue
		}

		log.Infof("schedule slots=%+v", slots)
		err := cluster.assignSlots(slots)
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
			unassignedSlots = append(unassignedSlots, common.SlotID(slot))
		}
	}

	return unassignedSlots
}

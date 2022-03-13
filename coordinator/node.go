package coordinator

import (
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
)

type Node interface {
	// GetId() common.NodeID
	GetSlots() []common.SlotID
	// HasSlot(id common.SlotID) bool
	// GetUsage() uint64
	// GetPercent() float32
	// GetCap() uint64
	// GetCpuPercent() float32

	AssignSlots(slots []*SlotInfo) error
}

type baseNode struct {
	info *proto.NodeInfo

	rwmutex             sync.RWMutex // guard fields below
	serveSlots          map[common.SlotID]*SlotInfo
	storageUsage        uint64
	storageUsagePercent float32
	cap                 uint64 // not guard, cap field is read-only
	cpuPercent          float32
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

func (node *baseNode) AssignSlots(slots []*SlotInfo) error {
	node.rwmutex.Lock()
	defer node.rwmutex.Unlock()

	for _, slot := range slots {
		node.serveSlots[slot.Id] = slot
	}

	return nil
}

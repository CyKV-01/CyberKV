package coordinator

import (
	"context"
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
)

type Compactor struct {
	rwmutex               sync.RWMutex
	activeMemTableCompact map[common.SlotID]struct{}
}

func NewCompactor() *Compactor {
	return &Compactor{
		rwmutex:               sync.RWMutex{},
		activeMemTableCompact: make(map[common.SlotID]struct{}),
	}
}

func (compactor *Compactor) Compact(slot common.SlotID, nodes []StorageNode) {
	compactor.rwmutex.RLock()
	if _, ok := compactor.activeMemTableCompact[slot]; ok {
		// The slot is compacting
		compactor.rwmutex.RUnlock()
		return
	}
	compactor.rwmutex.RUnlock()

	compactor.rwmutex.Lock()
	compactor.activeMemTableCompact[slot] = struct{}{}
	compactor.rwmutex.Unlock()

	ctx := context.Background()
	req := proto.CompactRequest{
		Leader:     nodes[0].info.Id,
		LeaderAddr: nodes[0].info.Addr,
		Slot:       int64(slot),
	}

	wg := sync.WaitGroup{}
	for _, node := range nodes {
		wg.Add(1)
		go func(node StorageNode) {
			defer wg.Done()
			node.StorageClient.Compact(ctx, &req)
		}(node)
	}
	wg.Wait()
}

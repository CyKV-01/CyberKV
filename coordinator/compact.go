package coordinator

import (
	"context"
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
)

type Compactor struct {
	rwmutex               sync.RWMutex
	activeMemTableCompact map[common.SlotID]struct{}
	versionSet            *VersionSet
}

func NewCompactor(versionSet *VersionSet) *Compactor {
	return &Compactor{
		rwmutex:               sync.RWMutex{},
		activeMemTableCompact: make(map[common.SlotID]struct{}),
		versionSet:            versionSet,
	}
}

func (compactor *Compactor) CompactMemTable(slot common.SlotID, nodes []*StorageNode) {
	compactor.rwmutex.RLock()
	if _, ok := compactor.activeMemTableCompact[slot]; ok {
		// The slot is compacting
		compactor.rwmutex.RUnlock()
		log.Info("the slot is already in compacting",
			zap.Int32("slot", slot))
		return
	}
	compactor.rwmutex.RUnlock()

	compactor.rwmutex.Lock()
	compactor.activeMemTableCompact[slot] = struct{}{}
	compactor.rwmutex.Unlock()

	ctx := context.Background()
	leader := nodes[0]
	req := proto.CompactMemTableRequest{
		Leader:     leader.info.Id,
		LeaderAddr: leader.info.Addr,
		Slot:       slot,
	}

	wg := sync.WaitGroup{}
	for _, node := range nodes {
		wg.Add(1)
		go func(node *StorageNode) {
			defer wg.Done()
			resp, err := node.CompactMemTable(ctx, &req)
			if err != nil {
				log.Error("failed to request storage node to compact",
					zap.String("node_id", node.info.Id),
					zap.Error(err))
				return
			}

			if node.info.Id == leader.info.Id {
				_, err = compactor.versionSet.Edit(resp.CreatedSstables, resp.DeletedSstables)
				if err != nil {
					log.Error("failed to update version set",
						zap.Error(err))
				}
			}
		}(node)
	}
	wg.Wait()

	compactor.rwmutex.Lock()
	delete(compactor.activeMemTableCompact, slot)
	compactor.rwmutex.Unlock()
}

package coordinator

import (
	"context"
	"sync"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
)

type Compactor struct {
	activeMemTableCompact *common.ConcurrentSet
	versionSet            *VersionSet
}

func NewCompactor(versionSet *VersionSet) *Compactor {
	return &Compactor{
		activeMemTableCompact: common.NewConcurrentSet(),
		versionSet:            versionSet,
	}
}

func (compactor *Compactor) CompactMemTable(slot common.SlotID, nodes []*StorageNode) {
	if compactor.activeMemTableCompact.Insert(slot) {
		// The slot is compacting
		log.Debug("the slot is already in compacting",
			zap.Int32("slot", slot))
		return
	}

	log.Info("start to compact slot",
		zap.Int32("slot", slot))

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
			c, _ := context.WithTimeout(ctx, 60*time.Second)
			resp, err := node.CompactMemTable(c, &req)
			if err != nil {
				log.Error("failed to request storage node to compact",
					zap.String("node_id", node.info.Id),
					zap.Error(err))
				return
			}

			if node.info.Id == leader.info.Id {
				log.Info("edit version set")
				version, err := compactor.versionSet.Edit(resp.CreatedSstables, resp.DeletedSstables)
				log.Info("edit version set done",
					zap.Uint64("new version", version.VersionId))

				if err != nil {
					log.Error("failed to update version set",
						zap.Error(err))
				}
			}
		}(node)
	}
	wg.Wait()

	log.Info("compaction done",
		zap.Int32("slot", slot))
	compactor.activeMemTableCompact.Delete(slot)
}

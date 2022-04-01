package storage

import (
	"sync"
	"sync/atomic"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
)

type SlotChan struct {
	idx int32
	ch  []chan *proto.KvData
}

func NewSlotChan() *SlotChan {
	channels := make([]chan *proto.KvData, common.DefaultReplicaNum)
	for i := range channels {
		channels[i] = make(chan *proto.KvData, 32)
	}
	return &SlotChan{
		idx: 0,
		ch:  channels,
	}
}

func (ch *SlotChan) AcquireChan() chan *proto.KvData {
	idx := atomic.AddInt32(&ch.idx, 1) - 1
	return ch.ch[idx]
}

func (ch *SlotChan) GetAllChan() []chan *proto.KvData {
	return ch.ch
}

type Compactor struct {
	rwmutex     sync.RWMutex
	compactChan map[common.SlotID]*SlotChan
}

func NewCompactor() *Compactor {
	return &Compactor{
		rwmutex:     sync.RWMutex{},
		compactChan: make(map[common.SlotID]*SlotChan),
	}
}

func (compactor *Compactor) PreCompact(slot common.SlotID) *SlotChan {
	ch := NewSlotChan()

	compactor.rwmutex.Lock()
	compactor.compactChan[slot] = ch
	compactor.rwmutex.Unlock()

	return ch
}

func (compactor *Compactor) GetChan(slot common.SlotID) *SlotChan {
	compactor.rwmutex.RLock()
	defer compactor.rwmutex.RUnlock()

	ch, ok := compactor.compactChan[slot]
	if !ok {
		return nil
	}

	return ch
}

func (compactor *Compactor) MergeChan(slotCh *SlotChan) chan *proto.KvData {
	mergeCh := make(chan *proto.KvData, 32)

	go func() {
		var (
			key             string
			currentKeyCount int
			ts              common.TimeStamp
			value           string
		)

		next := make(map[int]*proto.KvData)
		keyCount := make(map[string]int)

		for {
			for i, ch := range slotCh.ch {
				data, ok := next[i]
				if ok { // already in
					// log.Info("data is already in next map")
					continue
				}

				for {
					// log.Info("read data from channel")
					data, ok = <-ch
					if !ok {
						// log.Info("no more data in channel")
						break
					}
					// This key has been consumed
					if ts > 0 && data.Key == key {
						// log.Info("key has been consumed")
						continue
					}

					break
				}
				if !ok {
					continue
				}

				next[i] = data
			}

			// no more data
			if len(next) == 0 {
				// log.Info("no more data in all channels")
				break
			}

			// Find key
			currentKeyCount = 0
			for _, data := range next {
				count := keyCount[data.Key]
				count++

				if count > currentKeyCount {
					currentKeyCount = count
					key = data.Key
				}

				keyCount[data.Key] = count
			}
			delete(keyCount, key)

			// Find value and ts
			ts = 0
			for i, data := range next {
				if data.Key == key {
					delete(next, i)

					if data.Timestamp > ts {
						ts = data.Timestamp
						value = data.Value
					}
				}
			}

			// log.Info("insert kv into merge channel",
			// 	zap.String("key", key),
			// 	zap.String("value", value),
			// 	zap.Uint64("timestamp", ts))
			mergeCh <- &proto.KvData{
				Key:       key,
				Value:     value,
				Timestamp: ts,
			}
		}

		// log.Info("merge channels done")
		close(mergeCh)
	}()

	return mergeCh
}

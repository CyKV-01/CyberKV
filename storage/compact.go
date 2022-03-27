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
	return &SlotChan{
		idx: 0,
		ch: []chan *proto.KvData{
			make(chan *proto.KvData, 32),
			make(chan *proto.KvData, 32),
			make(chan *proto.KvData, 32),
		},
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
					continue
				}

				for {
					data, ok = <-ch
					if !ok {
						break
					}
					// This key has been consumed
					if ts > 0 && data.Key == key {
						continue
					}
				}
				if !ok {
					continue
				}
				
				next[i] = data
			}

			// no more data
			if len(next) == 0 {
				break
			}

			// Find key
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

			mergeCh <- &proto.KvData{
				Key:       value,
				Value:     value,
				Timestamp: ts,
			}
		}

		close(mergeCh)
	}()

	return mergeCh
}

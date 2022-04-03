package storage

import (
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
)

type SlotChan struct {
	mutex    sync.Mutex
	closed   bool
	globalCh chan chan *proto.KvData
}

func NewSlotChan() *SlotChan {
	return &SlotChan{
		mutex:    sync.Mutex{},
		closed:   false,
		globalCh: make(chan chan *proto.KvData, common.DefaultReplicaNum),
	}
}

func (ch *SlotChan) CreateDataChan() chan *proto.KvData {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	if ch.closed {
		return nil
	}

	dataCh := make(chan *proto.KvData, 32)
	ch.globalCh <- dataCh
	return dataCh
}

func (ch *SlotChan) Channels() []chan *proto.KvData {
	if !ch.closed {
		panic("Channels() can only be called after close slot channel")
	}

	channels := make([]chan *proto.KvData, 0, ch.Len())
	for ch := range ch.globalCh {
		channels = append(channels, ch)
	}

	return channels
}

func (ch *SlotChan) Len() int {
	return len(ch.globalCh)
}

func (ch *SlotChan) Close() {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	ch.closed = true
	close(ch.globalCh)
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

func (compactor *Compactor) MergeChan(slotCh *SlotChan, timestamp common.TimeStamp) chan *proto.KvData {
	mergeCh := make(chan *proto.KvData, 32)

	go func() {
		var (
			key string
			// currentKeyCount int
			ts    common.TimeStamp
			value string
		)

		next := make(map[int]*proto.KvData)
		// keyCount := make(map[string]int)

		channels := slotCh.Channels()
		log.Info("merge channels",
			zap.Int("channelNum", len(channels)))

		for {
			for i, ch := range channels {
				data, ok := next[i]
				if ok { // already in
					log.Debug("data is already in next map")
					continue
				}

				for {
					log.Debug("read data from channel")
					data, ok = <-ch
					if !ok {
						log.Debug("no more data in channel")
						break
					}
					// This key has been consumed
					if ts > 0 && data.Timestamp < timestamp && data.Key == key {
						log.Debug("key has been consumed",
							zap.String("key", data.Key),
							zap.Uint64("timestamp", data.Timestamp))
						continue
					}

					log.Debug("get data from channel",
						zap.Int("index", i),
						zap.String("key", data.Key),
						zap.String("value", data.Value))

					break
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
			ts = 0
			for _, data := range next {
				// count := keyCount[data.Key]
				// count++
				// keyCount[data.Key] = count

				if ts == 0 {
					key = data.Key
					value = data.Value
					ts = data.Timestamp
				} else if data.Key < key || data.Key == key && data.Timestamp > ts {
					key = data.Key
					value = data.Value
					ts = data.Timestamp
				}
			}

			for i, data := range next {
				if data.Key == key && data.Timestamp == ts {
					delete(next, i)
				}
			}

			// if currentKeyCount < common.DefaultReadQuorum {
			// 	log.Warn("failed to reach read quorum when merge channels",
			// 		zap.String("key", key),
			// 		zap.Int("count", currentKeyCount))
			// }
			// delete(keyCount, key)

			log.Debug("insert kv into merge channel",
				zap.String("key", key),
				zap.String("value", value),
				zap.Uint64("timestamp", ts))
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

package db

import (
	"sync"

	"github.com/yah01/CyberKV/common"
	. "github.com/yah01/CyberKV/common"
)

type MemTable[K Comparable[K], V any] struct {
	table *BTree[K, V]
}

func NewMemTable[K Comparable[K], V any]() *MemTable[K, V] {
	return &MemTable[K, V]{
		table: NewBTree[K, V](32),
	}
}

func (mem *MemTable[K, V]) Get(key K) *V {
	return mem.Get(key)
}

func (mem *MemTable[K, V]) Find(key K) (K, V, bool) {
	return mem.table.Find(key)
}

func (mem *MemTable[K, V]) Set(key K, value V) {
	mem.table.Insert(key, value)
}

func (mem *MemTable[K, V]) Range(fn func(key K, value V) bool) {
	mem.table.Range(fn)
}

type MemTableGroup[K Comparable[K], V any] [3]*MemTable[K, V]

func NewMemTableGroup[K Comparable[K], V any]() *MemTableGroup[K, V] {
	return &MemTableGroup[K, V]{NewMemTable[K, V]()}
}

func (group *MemTableGroup[K, V]) Mem() *MemTable[K, V] {
	return group[0]
}

func (group *MemTableGroup[K, V]) Imm() *MemTable[K, V] {
	return group[1]
}

func (group *MemTableGroup[K, V]) Fmem() *MemTable[K, V] {
	return group[2]
}

type SlotMemTable[K Comparable[K], V any] struct {
	rwmutex sync.RWMutex
	mem     map[SlotID]*MemTableGroup[K, V] // slot_id -> mem, imm, fmem

	sizeTableRwMutex sync.RWMutex
	memSize          map[SlotID]uint64
}

func NewSlotMemTable[K Comparable[K], V any]() *SlotMemTable[K, V] {
	return &SlotMemTable[K, V]{
		rwmutex: sync.RWMutex{},
		mem:     make(map[common.SlotID]*MemTableGroup[K, V]),

		sizeTableRwMutex: sync.RWMutex{},
		memSize:          make(map[int32]uint64),
	}
}

func (table *SlotMemTable[K, V]) GetMemTables(slot SlotID) *MemTableGroup[K, V] {
	table.rwmutex.RLock()
	defer table.rwmutex.RUnlock()

	group, ok := table.mem[slot]
	if !ok {
		return nil
	}

	return group
}

func (table *SlotMemTable[K, V]) GetMemTable(slot SlotID) *MemTable[K, V] {
	group, ok := table.mem[slot]
	if !ok {
		return nil
	}

	return group[0]
}

func (table *SlotMemTable[K, V]) CreateTables(slot SlotID) *MemTable[K, V] {
	group := NewMemTableGroup[K, V]()
	table.mem[slot] = group

	return group[0]
}

func (table *SlotMemTable[K, V]) Rotate(slot SlotID) *MemTableGroup[K, V] {
	table.rwmutex.Lock()
	defer table.rwmutex.Unlock()

	group := table.mem[slot]
	newGroup := *group
	for i := 2; i > 0; i-- {
		newGroup[i] = newGroup[i-1]
	}
	newGroup[0] = NewMemTable[K, V]()

	table.mem[slot] = &newGroup

	table.memSize = make(map[int32]uint64)

	return &newGroup
}

func (table *SlotMemTable[K, V]) Lock(slot SlotID) {
	table.rwmutex.Lock()
}

func (table *SlotMemTable[K, V]) Unlock(slot SlotID) {
	table.rwmutex.Unlock()
}

func (table *SlotMemTable[K, V]) AddSize(slot SlotID, size uint64) {
	// table.sizeTableRwMutex.Lock()
	// defer table.sizeTableRwMutex.Unlock()

	table.memSize[slot] += size
}

func (table *SlotMemTable[K, V]) GetSizeTable() map[SlotID]uint64 {
	table.rwmutex.RLock()
	defer table.rwmutex.RUnlock()

	clone := make(map[SlotID]uint64, len(table.memSize))
	for key, value := range table.memSize {
		clone[key] = value
	}

	return clone
}

// func (table *SlotMemTable[K, V]) Find(slot SlotID, key K) (K, *V) {
// 	table.rwmutex.RLock()
// 	defer table.rwmutex.RUnlock()

// 	mem, ok := table.mem[slot]
// 	if !ok {
// 		var resultKey K
// 		return resultKey, nil
// 	}

// 	return mem.Find(key)
// }

// func (table *SlotMemTable[K, V]) Set(slot SlotID, key K, value V) {
// 	table.rwmutex.Lock()
// 	defer table.rwmutex.Unlock()

// 	mem, ok := table.mem[slot]
// 	if !ok {
// 		mem = NewMemTable[K, V]()
// 		table.mem[slot] = mem
// 	}

// 	mem.Set(key, value)
// }

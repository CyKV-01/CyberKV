package db

import (
	"sync"

	. "github.com/yah01/CyberKV/common"
)

type MemTable[K Comparable[K], V any] struct {
	rwmutex sync.RWMutex
	table   *BTree[K, V]
}

func NewMemTable[K Comparable[K], V any]() *MemTable[K, V] {
	return &MemTable[K, V]{
		rwmutex: sync.RWMutex{},
		table:   NewBTree[K, V](32),
	}
}

func (mem *MemTable[K, V]) Get(key K) *V {
	mem.rwmutex.RLock()
	defer mem.rwmutex.RUnlock()

	return mem.Get(key)
}

func (mem *MemTable[K, V]) Find(key K) (K, *V) {
	mem.rwmutex.RLock()
	defer mem.rwmutex.RUnlock()

	return mem.table.Find(key)
}

func (mem *MemTable[K, V]) Set(key K, value V) {
	mem.rwmutex.Lock()
	defer mem.rwmutex.Unlock()

	mem.table.Insert(key, value)
}

type SlotMemTable[K Comparable[K], V any] struct {
	rwmutex   sync.RWMutex
	totalSize uint64
	mem       map[SlotID][]*MemTable[K, V] // slot_id -> mem, imm, fmem
}

func NewSlotMemTable[K Comparable[K], V any]() *SlotMemTable[K, V] {
	return &SlotMemTable[K, V]{
		rwmutex:   sync.RWMutex{},
		totalSize: 0,
		mem:       make(map[int16][]*MemTable[K, V]),
	}
}

func (table *SlotMemTable[K, V]) GetMemTables(slot SlotID) []*MemTable[K, V] {
	table.rwmutex.RLock()
	defer table.rwmutex.RUnlock()

	tables, ok := table.mem[slot]
	if !ok {
		return nil
	}

	return tables
}

func (table *SlotMemTable[K, V]) GetMemTable(slot SlotID) *MemTable[K, V] {
	tables, ok := table.mem[slot]
	if !ok {
		return nil
	}

	return tables[0]
}

func (table *SlotMemTable[K, V]) CreateTables(slot SlotID) *MemTable[K, V] {
	tables := make([]*MemTable[K, V], 3)
	tables[0] = NewMemTable[K, V]()

	table.mem[slot] = tables

	return tables[0]
}

func (table *SlotMemTable[K, V]) Lock(slot SlotID) {
	table.rwmutex.Lock()
}

func (table *SlotMemTable[K, V]) Unlock(slot SlotID) {
	table.rwmutex.Unlock()
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

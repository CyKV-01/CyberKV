package common

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	inner sync.Map
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		inner: sync.Map{},
	}
}

func (cmap *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	var v V
	vi, ok := cmap.inner.Load(key)
	if !ok {
		return v, false
	}

	v = vi.(V)
	return v, ok
}

func (cmap *ConcurrentMap[K, V]) GetOrInsert(key K, value V) (V, bool) {
	vi, exist := cmap.inner.LoadOrStore(key, value)
	if exist {
		return vi.(V), true
	}

	return value, false
}

func (cmap *ConcurrentMap[K, V]) Insert(key K, value V) {
	cmap.inner.Store(key, value)
}

func (cmap *ConcurrentMap[K, V]) Remove(key K) {
	cmap.inner.Delete(key)
}

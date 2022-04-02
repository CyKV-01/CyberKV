package common

import "sync"

type ConcurrentSet struct {
	set sync.Map
}

func NewConcurrentSet() *ConcurrentSet {
	return &ConcurrentSet{
		set: sync.Map{},
	}
}

func (set *ConcurrentSet) Insert(elem any) {
	set.set.Store(elem, struct{}{})
}

func (set *ConcurrentSet) Contain(elem any) bool {
	_, ok := set.set.Load(elem)
	return ok
}

func (set *ConcurrentSet) Delete(elem any) {
	set.set.Delete(elem)
}

type Set[T comparable] map[T]struct{}

func (set Set[T]) Insert(elem ...T) {
	for i := range elem {
		set[elem[i]] = struct{}{}
	}
}

func (set Set[T]) Contain(elem T) bool {
	_, ok := set[elem]
	return ok
}

func (set Set[T]) Delete(elem ...T) {
	for i := range elem {
		delete(set, elem[i])
	}
}

func (set Set[T]) Collect() []T {
	elems := make([]T, 0, len(set))
	for elem := range set {
		elems = append(elems, elem)
	}

	return elems
}

package common

import "sync"

type Set struct {
	set sync.Map
}

func NewSet() *Set {
	return &Set{
		set: sync.Map{},
	}
}

func (set *Set) Insert(elem any) {
	set.set.Store(elem, struct{}{})
}

func (set *Set) Contain(elem any) bool {
	_, ok := set.set.Load(elem)
	return ok
}

func (set *Set) Delete(elem any) {
	set.set.Delete(elem)
}

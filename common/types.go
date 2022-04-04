package common

type (
	SlotID    = int32
	NodeID    = uint64
	TimeStamp = uint64
	UniqueID  = uint64

	Comparator[T any] func(a, b T) int
)

type Comparable[T any] interface {
	Compare(other T) int
}

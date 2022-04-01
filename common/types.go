package common

type (
	SlotID    = int32
	NodeID    = string
	TimeStamp = uint64

	Comparator[T any] func(a, b T) int
)

type Comparable[T any] interface {
	Compare(other T) int
}

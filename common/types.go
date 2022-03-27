package common

type (
	SlotID    = uint16
	NodeID    = string
	TimeStamp = uint64

	Comparator[T any] func(a, b T) int
)

type Comparable[T any] interface {
	Compare(other T) int
}

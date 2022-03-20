package binary

import "unsafe"
import "golang.org/x/exp/constraints"

func Append[T constraints.Integer](buf *[]byte, data T) {
	size := int(unsafe.Sizeof(data))
	for i := 0; i < size; i++ {
		*buf = append(*buf, byte(data>>(i*8)))
	}
}

func Put[T constraints.Integer](buf []byte, data T) {
	size := int(unsafe.Sizeof(data))
	for i := 0; i < size; i++ {
		buf[i] = byte(data >> (i * 8))
	}
}

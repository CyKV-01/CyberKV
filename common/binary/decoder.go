package binary

import (
	"io"
	"unsafe"

	"golang.org/x/exp/constraints"
)

func Read[T constraints.Integer](reader io.ByteReader) (T, error) {
	var data T // zero

	size := int(unsafe.Sizeof(data))
	for i := 0; i < size; i++ {
		byt, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF && i < size-1 && i != 0 {
				err = io.ErrUnexpectedEOF
			}
			return data, err
		}

		data |= T(byt) << (i * 8)
	}

	return data, nil
}

func Get[T constraints.Integer](bytes []byte) T {
	var data T
	size := int(unsafe.Sizeof(data))
	for i := 0; i < size; i++ {
		data |= T(bytes[i]) << (i * 8)
	}

	return data
}

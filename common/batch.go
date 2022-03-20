package common

import "github.com/yah01/CyberKV/common/binary"

const (
	BatchHeaderSize = 8 // TimeStamp
)

type Batch struct {
	data []byte
}

func NewBatch() *Batch {
	return &Batch{
		data: make([]byte, BatchHeaderSize),
	}
}

func (batch *Batch) Put(key, value string) {
	AppendDataWithLength(&batch.data, key)
	AppendDataWithLength(&batch.data, value)
}

func (batch *Batch) SetSequence(ts TimeStamp) {
	binary.Put(batch.data, ts)
}

func (batch *Batch) Data() []byte {
	return batch.data
}

func AppendDataWithLength[T string | []byte](dst *[]byte, data T) {
	binary.Append(dst, uint16(len(data)))
	*dst = append(*dst, data...)
}

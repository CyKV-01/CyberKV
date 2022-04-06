package db

import (
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/binary"
	"github.com/yah01/CyberKV/common/log"
	"go.uber.org/zap"
)

const (
	BatchHeaderSize = 8 // TimeStamp
)

type Batch struct {
	data []byte
}

var (
	batchPool = sync.Pool{
		New: func() any {
			return &Batch{
				data: make([]byte, BatchHeaderSize, BatchHeaderSize+128),
			}
		},
	}
)

func NewBatch() *Batch {
	batch := batchPool.Get().(*Batch)
	batch.data = batch.data[:BatchHeaderSize]
	return batch
}

func NewBatchFromBytes(bytes []byte) *Batch {
	return &Batch{
		data: bytes,
	}
}

func (batch *Batch) Close() {
	batchPool.Put(batch)
}

func (batch *Batch) Put(key, value string) {
	binary.Append(&batch.data, common.SetValueType)
	AppendDataWithLength(&batch.data, key)
	AppendDataWithLength(&batch.data, value)
}

func (batch *Batch) SetSequence(ts common.TimeStamp) {
	binary.Put(batch.data, ts)
}

func (batch *Batch) Data() []byte {
	return batch.data
}

func (batch *Batch) GetSequence() common.TimeStamp {
	return binary.Get[common.TimeStamp](batch.data)
}

func (batch *Batch) GetKvs() (keys []string, values []string, types []common.ValueType) {
	var (
		valueType common.ValueType
		keyLen    uint16
		valueLen  uint16
		key       string
		value     string
		buf       = make([]byte, 32)
	)

	data := batch.data[BatchHeaderSize:]
	offset := uint64(0)
	for offset < uint64(len(data)) {
		valueType = binary.Get[common.ValueType](data[offset:])
		offset += 1

		keyLen = binary.Get[uint16](data[offset:])
		offset += 2

		if len(buf) < int(keyLen) {
			buf = make([]byte, keyLen)
		}

		copy(buf[:keyLen], data[offset:])
		offset += uint64(keyLen)
		key = string(buf[:keyLen])

		valueLen = binary.Get[uint16](data[offset:])
		offset += 2

		if len(buf) < int(valueLen) {
			buf = make([]byte, valueLen)
		}

		copy(buf[:valueLen], data[offset:])
		offset += uint64(valueLen)
		value = string(buf[:valueLen])

		keys = append(keys, key)
		values = append(values, value)
		types = append(types, valueType)

		log.Debug("parse key-value from batch",
			zap.String("key", key),
			zap.String("value", value),
			zap.Uint8("valueType", uint8(valueType)))
	}

	return
}

func AppendDataWithLength[T string | []byte](dst *[]byte, data T) {
	binary.Append(dst, uint16(len(data)))
	*dst = append(*dst, data...)
}

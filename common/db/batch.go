package db

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/binary"
)

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

func NewBatchFromBytes(bytes []byte) *Batch {
	return &Batch{
		data: bytes,
	}
}

func (batch *Batch) Put(key, value string) {
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

func (batch *Batch) GetKvs() (keys []string, values []string, err error) {
	var (
		keyLen   uint16
		valueLen uint16
		key      string
		value    string
		buf      = make([]byte, 32)
		n        int
	)

	reader := bufio.NewReader(bytes.NewReader(batch.data[BatchHeaderSize:]))
	for {
		keyLen, err = binary.Read[uint16](reader)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		if len(buf) < int(keyLen) {
			buf = make([]byte, keyLen)
		}
		n, err = io.ReadFull(reader, buf[:keyLen])
		if err != nil {
			return
		}
		if n != int(keyLen) {
			return nil, nil, fmt.Errorf("failed to read key, keyLen=%v n=%v", keyLen, n)
		}
		key = string(buf[:keyLen])

		valueLen, err = binary.Read[uint16](reader)
		if err != nil {
			return
		}
		if len(buf) < int(valueLen) {
			buf = make([]byte, valueLen)
		}
		n, err = io.ReadFull(reader, buf[:valueLen])
		if err != nil && err != io.EOF {
			return
		}
		if n != int(valueLen) {
			return nil, nil, fmt.Errorf("failed to read value, valueLen=%v n=%v err=%v", valueLen, n, err)
		}
		value = string(buf[:valueLen])

		keys = append(keys, key)
		values = append(values, value)

		if err == io.EOF {
			err = nil
			break
		}
	}

	return
}

func AppendDataWithLength[T string | []byte](dst *[]byte, data T) {
	binary.Append(dst, uint16(len(data)))
	*dst = append(*dst, data...)
}

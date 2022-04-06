package db

import (
	"io"
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/binary"
)

const (
	Full RecordType = iota + 1
	First
	Middle
	Last

	RecordHeaderSize = 4 + 2 + 1
)

type RecordType uint8

type Record struct {
	Checksum uint32
	Length   uint16
	Type     RecordType
	Data     []byte
}

func NewRecord(data []byte) *Record {
	record := new(Record)
	record.Checksum = common.CalcChecksum(data)
	record.Length = uint16(len(data))
	record.Type = Full
	record.Data = data

	return record
}

func NewRecordFromBytes(data []byte) (*Record, uint64) {
	var (
		record Record
		size   uint64
	)

	if len(data) < RecordHeaderSize {
		return nil, 0
	}

	record.Checksum = binary.Get[uint32](data)
	size += 4
	record.Length = binary.Get[uint16](data[size:])
	size += 2
	record.Type = binary.Get[RecordType](data[size:])
	size += 1
	if len(data[size:]) < int(record.Length) {
		return nil, 0
	}

	record.Data = make([]byte, record.Length)
	copy(record.Data, data[size:size+uint64(record.Length)])
	size += uint64(record.Length)

	return &record, size
}

func (record *Record) BuildBytes() []byte {
	totalSize := RecordHeaderSize + len(record.Data)
	buffer := make([]byte, 0, totalSize)

	binary.Append(&buffer, record.Checksum)
	binary.Append(&buffer, record.Length)
	binary.Append(&buffer, record.Type)
	buffer = append(buffer, record.Data...)

	return buffer
}

var (
	metaBufferPool = sync.Pool{
		New: func() any {
			return make([]byte, 0, RecordHeaderSize)
		},
	}
)

func (record *Record) WriteTo(writer io.Writer) error {
	metaBuffer := metaBufferPool.Get().([]byte)[:0]
	binary.Append(&metaBuffer, record.Checksum)
	binary.Append(&metaBuffer, record.Length)
	binary.Append(&metaBuffer, record.Type)

	_, err := writer.Write(metaBuffer)
	metaBufferPool.Put(metaBuffer)
	if err != nil {
		return err
	}

	_, err = writer.Write(record.Data)
	return err
}

package db

import (
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

type RecordType = uint8

type Record struct {
	Checksum uint32
	Length   uint16
	Type     uint8
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

func (record *Record) BuildBytes() []byte {
	totalSize := RecordHeaderSize + len(record.Data)
	buffer := make([]byte, 0, totalSize)

	binary.Append(&buffer, record.Checksum)
	binary.Append(&buffer, record.Length)
	binary.Append(&buffer, record.Type)
	buffer = append(buffer, record.Data...)

	return buffer
}

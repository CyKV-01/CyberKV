package db

import (
	"bufio"
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

func NewRecordFromByteReader(reader *bufio.Reader) (*Record, error) {
	var (
		record Record
		err    error
	)

	record.Checksum, err = binary.Read[uint32](reader)
	if err != nil {
		return nil, err
	}

	record.Length, err = binary.Read[uint16](reader)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}

	record.Type, err = binary.Read[RecordType](reader)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return nil, err
	}

	record.Data = make([]byte, record.Length)
	_, err = io.ReadFull(reader, record.Data)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &record, err
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

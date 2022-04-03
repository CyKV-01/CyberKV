package db

import (
	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/binary"
	"github.com/yah01/CyberKV/proto"
)

type KvData struct {
	Timestamp common.TimeStamp
	KeyLen    uint16
	ValueLen  uint16
	Key       string
	Value     string
}

func (kv *KvData) Len() uint64 {
	return 8 + 2 + 2 + uint64(kv.KeyLen+kv.ValueLen)
}

type SSTable struct {
	data     []byte
	Level    int
	Path     string
	FirstKey string
	LastKey  string

	index *proto.Index
}

func NewSSTable(data []byte, level int, path, firstKey, lastKey string, index *proto.Index) *SSTable {
	return &SSTable{
		data:     data,
		Level:    level,
		Path:     path,
		FirstKey: firstKey,
		LastKey:  lastKey,
		index:    index,
	}
}

func NewSSTableFromDataCh(dataCh chan *proto.KvData) *SSTable {
	var (
		sstableBytes      = make([]byte, 0, 16<<20)
		index             proto.Index
		blockSize         uint64
		oldOffset         uint64
		newOffset         uint64
		shouldDivideBlock = true
		first             *proto.KvData
		last              *proto.KvData
	)

	for last = range dataCh {
		if first == nil {
			first = last
		}
		oldOffset = uint64(len(sstableBytes))

		binary.Append(&sstableBytes, last.Timestamp)
		binary.Append(&sstableBytes, uint16(len(last.Key)))
		binary.Append(&sstableBytes, uint16(len(last.Value)))
		sstableBytes = append(sstableBytes, last.Key...)
		sstableBytes = append(sstableBytes, last.Value...)

		newOffset = uint64(len(sstableBytes))

		if shouldDivideBlock {
			shouldDivideBlock = false
			blockSize = 0
			index.BlockHandles = append(index.BlockHandles, &proto.BlockHandle{
				Key:    last.Key,
				Offset: uint64(oldOffset),
			})
		}
		blockSize += newOffset - oldOffset

		if blockSize >= common.SSTableBlockSize {
			shouldDivideBlock = true
		}
	}

	if first == nil {
		return nil
	}

	return NewSSTable(sstableBytes, 0, "", first.Key, last.Key, &index)
}

func (table *SSTable) Empty() bool {
	return len(table.data) == 0
}

func (table *SSTable) Len() int {
	return len(table.data)
}
func (table *SSTable) GetData() []byte {
	return table.data
}

func (table *SSTable) GetIndex() *proto.Index {
	return table.index
}

func (table *SSTable) Get(key string, timestamp common.TimeStamp) (string, bool) {
	var (
		l   = 0
		r   = len(table.index.BlockHandles) - 1
		mid int
	)

	for l < r {
		mid = (l+r)>>1 + 1
		if table.index.BlockHandles[mid].Key <= key {
			l = mid
		} else {
			r = mid - 1
		}
	}

	offset := table.index.BlockHandles[r].Offset
	endOffset := uint64(len(table.data))
	if r+1 < len(table.index.BlockHandles) {
		endOffset = table.index.BlockHandles[r+1].Offset
	}

	for offset < endOffset {
		kv := table.GetByOffset(offset)
		if kv.Key == key && kv.Timestamp <= timestamp {
			return kv.Value, true
		}

		offset += kv.Len()
	}

	return "", false
}

func (table *SSTable) GetByOffset(offset uint64) KvData {
	var kv KvData

	kv.Timestamp = binary.Get[common.TimeStamp](table.data[offset:])
	offset += 8

	kv.KeyLen = binary.Get[uint16](table.data[offset:])
	offset += 2

	kv.ValueLen = binary.Get[uint16](table.data[offset:])
	offset += 2

	kv.Key = string(table.data[offset : offset+uint64(kv.KeyLen)])
	offset += uint64(kv.KeyLen)

	kv.Value = string(table.data[offset : offset+uint64(kv.KeyLen)])

	return kv
}

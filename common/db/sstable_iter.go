package db

type SSTableIter struct {
	table         *SSTable
	currentOffset uint64
	value         KvData
}

func NewSSTableIter(sstable *SSTable) *SSTableIter {
	return &SSTableIter{
		table:         sstable,
		currentOffset: 0,
	}
}

func (iter *SSTableIter) Next() {
	iter.currentOffset += iter.value.Len()
}

func (iter *SSTableIter) Value() KvData {
	iter.value = iter.table.GetByOffset(iter.currentOffset)
	return iter.value
}

func (iter *SSTableIter) Valid() bool {
	return iter.currentOffset < uint64(iter.table.Len())
}

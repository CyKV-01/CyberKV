package common

const (
	WalCompactThreshold  = 4 * 1024 * 200
	MaxMemTableTotalSize = 32 * 1024 * 1024
	SSTableBlockSize     = 4 * 1024

	SlotNum = 1

	ServicePrefix = "services"
	SlotPrefix    = "slots"
	SSTableIdKey  = "sstable_id"
	TimestampKey  = "timestamp"
	VersionSetKey = "version_set"
	VersionPrefix = "version"

	DefaultTTL = 10

	DefaultReplicaNum  = 3
	DefaultReadQuorum  = 2
	DefaultWriteQuorum = 2
	MaxLevel           = 8
)

const (
	SetValueType ValueType = iota + 1
	DeleteValueType
)

func init() {
	if DefaultReadQuorum+DefaultWriteQuorum <= DefaultReplicaNum {
		panic("read quorum + write quorum must be greater than replica number")
	}
}

package common

const (
	WalCompactThreshold  = 4 * 1024 * 100
	MaxMemTableTotalSize = 32 * 1024 * 1024
	SSTableBlockSize     = 4 * 1024

	SlotNum = 1

	ServicePrefix = "services"
	SlotPrefix    = "slots"
	SSTableIdKey  = "sstable_id"
	VersionSetKey = "version_set"
	VersionPrefix = "version"

	DefaultTTL = 10

	DefaultReplicaNum  = 3
	DefaultReadQuorum  = 2
	DefaultWriteQuorum = 2
	MaxLevel           = 8
)

func init() {
	if DefaultReadQuorum+DefaultWriteQuorum <= DefaultReplicaNum {
		panic("read quorum + write quorum must be greater than replica number")
	}
}

package common

const (
	WalCompactThreshold          = 4 * 1024 * 1024
	MaxMemTableTotalSize         = 32 * 1024 * 1024
	Level0CompactThreshold       = 16 * 1024 * 1024
	Level0NumberCompactThreshold = 4

	SlotNum = 32

	ServicePrefix = "services"
	SlotPrefix    = "slots"

	DefaultTTL = 10

	DefaultReplicaNum  = 1
	DefaultReadQuorum  = 1
	DefaultWriteQuorum = 1
)

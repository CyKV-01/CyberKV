package common

import "hash/crc32"

type SlotID int16
type NodeID string

const (
	SlotNum = 32

	ServicePrefix = "services"
	DefaultTTL    = 30
)

func CalcSlotID(key string) SlotID {
	crc := crc32.Checksum([]byte(key), crc32.MakeTable(crc32.Castagnoli))
	return SlotID(crc % SlotNum)
}

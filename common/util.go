package common

import (
	"fmt"
	"hash/crc32"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
)

const DataDir = "data"

func LogName(slot SlotID, id int64) string {
	return fmt.Sprintf("%v_%v.log", slot, id)
}

func LogPath(slot SlotID, id int64) string {
	return path.Join(DataDir, LogName(slot, id))
}

func CalcSlotID(key string) SlotID {
	crc := crc32.Checksum([]byte(key), crc32.MakeTable(crc32.Castagnoli))
	return SlotID(crc % SlotNum)
}

func CalcChecksum(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
}

var uniqueId int64

func GenerateUniqueId() int64 {
	return atomic.AddInt64(&uniqueId, 1)
}

func GetSlotIdByEtcdPath(path string) SlotID {
	slotStr := strings.Split(path, "/")[1]
	slotId, err := strconv.ParseUint(slotStr, 10, 64)
	if err != nil {
		panic(err)
	}

	return SlotID(slotId)
}

func SliceContains[T comparable](slice []T, target T) bool {
	for i := range slice {
		if slice[i] == target {
			return true
		}
	}

	return false
}

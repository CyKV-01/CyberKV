package common

import (
	"fmt"
	"hash/crc32"
	"path"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/yah01/CyberKV/proto"
)

const DataDir = "data"

func LogName(slot SlotID, id int64) string {
	return fmt.Sprintf("%v_%v.log", slot, id)
}

func LogPath(slot SlotID, nodeID string, id int64) string {
	return path.Join(DataDir, nodeID, LogName(slot, id))
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

func IsOk(status *proto.Status) bool {
	return status == nil || status.ErrCode == proto.ErrorCode_Ok
}

// Search for the first element which is not less than target
func Search[S ~[]T, T Comparable[T]](slice S, target T) int {
	if len(slice) == 0 {
		return 0
	}

	var (
		l   = 0
		r   = len(slice) - 1
		mid int
	)

	for l < r {
		mid = (l + r) >> 1
		if slice[mid].Compare(target) < 0 {
			l = mid + 1
		} else {
			r = mid
		}
	}

	if slice[r].Compare(target) < 0 {
		return len(slice)
	}

	return r
}

func SSTableDataLevelPrefix(level int) string {
	return fmt.Sprintf("sstable/%d", level)
}

func SSTableDataPath(level int, id uint64) string {
	return fmt.Sprintf("%s/%d.cdb", SSTableDataLevelPrefix(level), id)
}

func SSTableIndexLevelPrefix(level int) string {
	return fmt.Sprintf("index/%d", level)
}

func SSTableIndexPath(level int, id uint64) string {
	return fmt.Sprintf("%s/%d.idx", SSTableIndexLevelPrefix(level), id)
}

func GetSSTableIndexPath(dataPath string) string {
	parts := strings.Split(dataPath, "/")
	parts[0] = "index"
	length := len(parts[2])
	parts[2] = parts[2][:length-3] + "idx"

	return strings.Join(parts, "/")
}

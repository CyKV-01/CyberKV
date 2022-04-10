package storage

import (
	"encoding/json"
	"os"
	"path"
	"sort"
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"go.uber.org/zap"
)

type Version struct {
	rwmutex sync.RWMutex             `json:"-"`
	Wals    map[common.SlotID]string `json:"wals"` // slotID -> log_name
}

func NewVersion() *Version {
	return &Version{
		Wals: make(map[int32]string),
	}
}

func RecoverVersion(node *StorageNode, logDir string) {
	log.Info("recover version...")

	entries, err := os.ReadDir(logDir)
	if err != nil {
		panic(err)
	}

	logs := make([]string, 0, len(node.version.Wals))
	for i := range entries {
		log.Info("recover version",
			zap.String("logName", entries[i].Name()))
		slotID, logID := common.ParseLogName(entries[i].Name())
		if path, ok := node.version.Wals[slotID]; ok {
			_, versionLogID := common.ParseLogName(path)
			if logID >= versionLogID {
				logs = append(logs, entries[i].Name())
			}
		} else {
			logs = append(logs, entries[i].Name())
		}
	}

	sort.Slice(logs, func(i, j int) bool {
		slotID0, logID0 := common.ParseLogName(logs[i])
		slotID1, logID1 := common.ParseLogName(logs[j])

		return slotID0 < slotID1 || slotID0 == slotID1 && logID0 < logID1
	})

	for i := 0; i < len(logs); i++ {
		slotID, _ := common.ParseLogName(logs[i])
		j := 0
		for j+1 < len(logs) {
			jSlotID, _ := common.ParseLogName(logs[j+1])
			if slotID != jSlotID {
				break
			}
			j++
		}

		mem := node.mem.CreateTables(slotID)
		for ; i <= j; i++ {
			log.Info("recover wal",
				zap.String("log", logs[i]))
			_, logID := common.ParseLogName(logs[i])
			reader, err := NewLogReader(path.Join(logDir, logs[i]))
			if err != nil {
				panic(err)
			}

			for {
				record := reader.NextRecord()
				if record == nil {
					break
				}

				batch := db.NewBatchFromBytes(record.Data)
				ts := batch.GetSequence()
				keys, values, types := batch.GetKvs()
				for i := range keys {
					if types[i] == common.SetValueType {
						log.Info("recover value",
							zap.String("key", keys[i]),
							zap.String("value", values[i]),
							zap.Uint64("timestamp", ts))
						mem.Set(db.NewInternalKey(keys[i], ts), values[i])
					}
				}
			}

			if logID > node.logID {
				node.logID = logID
			}
		}
	}

	for i := 0; i < len(logs); i++ {
		slotID, _ := common.ParseLogName(logs[i])
		node.logID++
		node.wals.Insert(slotID,
			NewLogWriter(slotID, node.Info.Id, node.logID))
	}

	log.Info("recover version done",
		zap.Uint64("logID", node.logID))
}

func (version *Version) Set(slotID common.SlotID, logName string) {
	version.rwmutex.Lock()
	defer version.rwmutex.Unlock()

	version.Wals[slotID] = logName

	versionBytes, err := json.Marshal(version)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("version.json", versionBytes, 0666)
	if err != nil {
		panic(err)
	}
}

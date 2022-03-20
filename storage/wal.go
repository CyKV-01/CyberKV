package storage

import (
	"os"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
)

type LogWriter struct {
	file *os.File
}

func NewLogWritter(slot common.SlotID) *LogWriter {
	id := common.GenerateUniqueId()
	file, err := os.OpenFile(common.LogPath(slot, id), os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("failed to open log file, err=%v", err)
		panic(err)
	}

	return &LogWriter{
		file: file,
	}
}

func (writer *LogWriter) Append(batch *common.Batch) error {
	record := common.NewRecord(batch.Data())

	_, err := writer.file.Write(record.BuildBytes())
	if err != nil {
		return err
	}

	err = writer.file.Sync()
	return err
}

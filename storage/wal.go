package storage

import (
	"os"
	"path"

	"github.com/yah01/CyberKV/common"
)

type LogWriter struct {
	file *os.File
}

func NewLogWritter(slot common.SlotID) *LogWriter {
	id := common.GenerateUniqueId()
	logPath := common.LogPath(slot, id)

	err := os.MkdirAll(path.Dir(logPath), 0777)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
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

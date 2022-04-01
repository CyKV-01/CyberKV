package storage

import (
	"bufio"
	"io"
	"io/fs"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"go.uber.org/zap"
)

type LogWriter struct {
	file *os.File
}

func NewLogWriter(slot common.SlotID, nodeID string) *LogWriter {
	id := common.GenerateUniqueId()
	logPath := common.LogPath(slot, nodeID, id)

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

func (writer *LogWriter) Append(batch *db.Batch) error {
	record := db.NewRecord(batch.Data())

	_, err := writer.file.Write(record.BuildBytes())
	if err != nil {
		return err
	}

	err = writer.file.Sync()
	return err
}

type LogReader struct {
	slot     common.SlotID
	file     *os.File
	curLogId int64

	buf *bufio.Reader
}

func NewLogReader(slot common.SlotID) (*LogReader, error) {

	file, logId, err := OpenNextLog(slot, -1)
	if err != nil {
		return nil, err
	}

	return &LogReader{
		slot:     slot,
		file:     file,
		curLogId: logId,
		buf:      bufio.NewReader(file),
	}, nil
}

func OpenNextLog(slot common.SlotID, curLogId int64) (file *os.File, newLogId int64, err error) {
	newLogId = math.MaxInt64

	err = filepath.WalkDir(common.DataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".log") {
			return nil
		}

		var (
			logSlot int64
			logId   int64
		)

		parts := strings.Split(d.Name(), "_")
		logSlot, err = strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			return err
		}
		parts[1] = strings.TrimSuffix(parts[1], ".log")
		logId, err = strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return err
		}

		if logSlot == int64(slot) && logId > curLogId && logId < newLogId {
			newLogId = logId
		}

		return nil
	})

	if err != nil {
		return
	}

	if newLogId == math.MaxInt64 {
		err = io.EOF
		return
	}

	logPath := ""
	file, err = os.OpenFile(logPath, os.O_RDONLY, 0666)

	return
}

func (reader *LogReader) NextRecord() (*db.Record, error) {
	if reader.file == nil {
		return nil, io.EOF
	}

	log.Info("reading record")
	record, err := db.NewRecordFromByteReader(reader.buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if err == io.EOF {
		// switch to next log, and re-read
		log.Info("switch to next log")
		reader.file, reader.curLogId, err = OpenNextLog(reader.slot, reader.curLogId)
		if err != nil {
			log.Info("no more log",
				zap.Int32("slot", reader.slot),
				zap.Int64("log_id", reader.curLogId))
			return nil, err
		}

		reader.buf = bufio.NewReader(reader.file)

		if record == nil {
			log.Info("re-read")
			record, err = reader.NextRecord()
		}
	}

	log.Infof("got record=%+v", record)
	return record, err
}

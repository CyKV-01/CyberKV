package storage

import (
	"io"
	"io/fs"
	"io/ioutil"
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
	path string
}

func NewLogWriter(slot common.SlotID, nodeID common.NodeID, logID common.UniqueID) *LogWriter {
	nodeIdStr := strconv.FormatUint(nodeID, 10)
	logPath := common.LogPath(slot, nodeIdStr, logID)

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
		path: logPath,
	}
}

func (writer *LogWriter) Append(batch *db.Batch) error {
	record := db.NewRecord(batch.Data())

	err := record.WriteTo(writer.file)
	if err != nil {
		return err
	}

	err = writer.file.Sync()
	return err
}

func (writer *LogWriter) Name() string {
	return path.Base(writer.path)
}

func (writer *LogWriter) Path() string {
	return writer.path
}

func (writer *LogWriter) Close() {
	writer.file.Close()
}

type LogReader struct {
	data   []byte
	offset uint64
}

func NewLogReader(path string) (*LogReader, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &LogReader{
		data:   data,
		offset: 0,
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

func (reader *LogReader) NextRecord() *db.Record {
	log.Info("reading record")
	record, size := db.NewRecordFromBytes(reader.data[reader.offset:])
	reader.offset += size

	log.Debug("got record",
		zap.Any("record", record))
	return record
}

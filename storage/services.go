package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (node *StorageNode) Get(ctx context.Context, request *proto.ReadRequest) (*proto.ReadResponse, error) {
	slot := common.CalcSlotID(request.Key)

	log.Info("opening logs...")
	reader, err := NewLogReader(slot)
	if err != nil {
		log.Error("failed to open log reader",
			zap.Int16("slot", slot),
			zap.String("key", request.Key))
		return &proto.ReadResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_IoError,
				ErrMessage: fmt.Sprintf("failed to read WAL, err=%v", err),
			}}, nil
	}

	log.Info("reading records...")
	ts := common.TimeStamp(0)
	value := ""
	for {
		record, err := reader.NextRecord()
		if err != nil && err != io.EOF {
			return &proto.ReadResponse{
				Status: &proto.Status{
					ErrCode:    proto.ErrorCode_IoError,
					ErrMessage: fmt.Sprintf("failed to decode record, err=%v", err),
				}}, nil
		}
		if record == nil {
			break
		}

		log.Info("reading batch")
		batch := db.NewBatchFromBytes(record.Data)

		keys, values, err := batch.GetKvs()
		if err != nil {
			return &proto.ReadResponse{
				Status: &proto.Status{
					ErrCode:    proto.ErrorCode_IoError,
					ErrMessage: fmt.Sprintf("failed to decode batch, err=%v", err),
				}}, nil
		}

		recordTs := batch.GetSequence()
		log.Info("batch info",
			zap.Uint64("batch_ts", recordTs))

		for i := range keys {
			if keys[i] == request.Key && recordTs > ts {
				ts = recordTs
				value = values[i]
			}
		}

		log.Info("updated from log",
			zap.Uint64("ts", ts),
			zap.String("value", value))
	}

	return &proto.ReadResponse{
		Ts:    ts,
		Value: value,
	}, nil
}

func (node *StorageNode) Set(ctx context.Context, request *proto.WriteRequest) (*proto.WriteResponse, error) {
	slot := common.CalcSlotID(request.Key)

	node.walMutex.Lock()
	wal, ok := node.wals[slot]
	if !ok {
		wal = NewLogWriter(slot)
		node.wals[slot] = wal
	}
	node.walMutex.Unlock()

	batch := db.NewBatch()
	batch.SetSequence(request.Ts)
	batch.Put(request.Key, request.Value)

	err := wal.Append(batch)
	if err != nil {
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_IoError,
				ErrMessage: fmt.Sprintf("failed to write WAL, err=%v", err),
			}}, nil
	}

	return &proto.WriteResponse{}, nil
}

func (node *StorageNode) Remove(ctx context.Context, request *proto.WriteRequest) (*proto.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

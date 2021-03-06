package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/common/wait"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KeyValue service

func (node *StorageNode) Get(ctx context.Context, request *proto.ReadRequest) (*proto.ReadResponse, error) {
	log.Debug("receive Get request",
		zap.String("key", request.Key),
		zap.Uint64("timestamp", request.Ts))

	slot := common.CalcSlotID(request.Key)

	mem, imm := node.mem.GetMemTables(slot)
	if mem == nil {
		return &proto.ReadResponse{
			Status: &proto.Status{
				ErrCode: proto.ErrorCode_KeyNotFound,
			},
		}, nil
	}

	var err error
	internalKey := db.NewInternalKey(request.Key, request.Ts)

	key, value, ok := mem.Find(internalKey)
	if (!ok || key.UserKey() != request.Key) && imm != nil {
		key, value, ok = imm.Find(internalKey)
		if !ok || key.UserKey() != request.Key {
			log.Info("key not found in memtables, will try to read it from sstables",
				zap.String("key", request.Key))
			value, ok, err = node.tableMgr.Get(ctx, internalKey)
			if err != nil {
				return &proto.ReadResponse{
					Status: &proto.Status{
						ErrCode:    proto.ErrorCode_IoError,
						ErrMessage: err.Error(),
					},
				}, nil
			}
		}
	}

	if !ok {
		return &proto.ReadResponse{
			Status: &proto.Status{
				ErrCode: proto.ErrorCode_KeyNotFound,
			},
		}, nil
	}

	log.Debug("respond to Get request",
		zap.String("key", request.Key),
		zap.String("value", value),
		zap.Uint64("timestamp", key.GetTimeStamp()))

	return &proto.ReadResponse{
		Value: value,
		Ts:    key.GetTimeStamp(),
	}, nil
}

func (node *StorageNode) Set(ctx context.Context, request *proto.WriteRequest) (*proto.WriteResponse, error) {
	slot := common.CalcSlotID(request.Key)

	wal, ok := node.wals.Get(slot)
	if !ok {
		log.Error("invalid slot",
			zap.Int32("slot", slot))
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_InvalidArgument,
				ErrMessage: "invalid slot",
			},
		}, nil
	}

	batch := db.NewBatch()
	batch.SetSequence(request.Ts)
	batch.Put(request.Key, request.Value)

	err := wal.Append(batch)
	batch.Close()
	if err != nil {
		log.Error("failed to write WAL",
			zap.Error(err))
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_IoError,
				ErrMessage: fmt.Sprintf("failed to write WAL, err=%v", err),
			}}, nil
	}

	internalKey := db.NewInternalKey(request.Key, request.Ts)
	node.mem.Lock(slot)
	mem := node.mem.GetMemTable(slot)
	if mem == nil {
		mem = node.mem.CreateTables(slot)
	}
	mem.Set(internalKey, request.Value)
	node.mem.AddSize(slot, uint64(len(request.Key)+len(request.Value)))
	node.mem.Unlock(slot)

	return &proto.WriteResponse{}, nil
}

func (node *StorageNode) Remove(ctx context.Context, request *proto.WriteRequest) (*proto.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

// Storage service
// Coordinator to storage nodes
func (node *StorageNode) CompactMemTable(ctx context.Context, request *proto.CompactMemTableRequest) (*proto.CompactMemTableResponse, error) {
	if request.Leader == node.Info.Id {
		return node.CompactMemTableAsLeader(ctx, request)
	} else {
		return node.CompactMemTableAsFollower(ctx, request)
	}
}

// Storage nodes to leader storage node
func (node *StorageNode) PushMemTable(ctx context.Context, request *proto.PushMemTableRequest) (*proto.PushMemTableResponse, error) {
	var slotCh *SlotChan
	ok := wait.WaitForCondition(10*time.Second, func() bool {
		slotCh = node.compactor.GetChan(request.Slot)
		return slotCh != nil
	})

	if !ok {
		log.Error("failed to get slot channel",
			zap.Int32("slot", request.Slot))
		return &proto.PushMemTableResponse{}, nil
	}

	log.Info("compaction follower create data channel",
		zap.Int32("slot", request.Slot))
	ch := slotCh.CreateDataChan()
	if ch == nil {
		log.Info("no need for more follower for compaction",
			zap.Int32("slot", request.Slot))
		return &proto.PushMemTableResponse{}, nil
	}

	for _, data := range request.Data {
		ch <- data
	}

	close(ch)

	return &proto.PushMemTableResponse{}, nil
}

func (node *StorageNode) AssignSlot(ctx context.Context, request *proto.AssignSlotRequest) (*proto.AssignSlotResponse, error) {
	log.Info("recv AssignSlot request",
		zap.Int32("slot", request.SlotID))

	wal := NewLogWriter(request.SlotID, node.Info.Id, node.nextLogID())
	old, ok := node.wals.GetOrInsert(request.SlotID, wal)
	if ok {
		log.Warn("The slot is assigned")
		wal.Close()
		os.Remove(old.Path())
		return &proto.AssignSlotResponse{}, nil
	}

	return &proto.AssignSlotResponse{}, nil
}

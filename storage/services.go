package storage

import (
	"context"
	"fmt"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KeyValue service

func (node *StorageNode) Get(ctx context.Context, request *proto.ReadRequest) (*proto.ReadResponse, error) {
	log.Info("get request",
		zap.String("key", request.Key),
		zap.Uint64("timestamp", request.Ts))

	slot := common.CalcSlotID(request.Key)

	tables := node.mem.GetMemTables(slot)
	if tables == nil {
		return &proto.ReadResponse{
			Status: &proto.Status{
				ErrCode: proto.ErrorCode_KeyNotFound,
			},
		}, nil
	}

	mem, imm, fmem := tables[0], tables[1], tables[2]
	internalKey := db.NewInternalKey(request.Key, request.Ts)
	key, value, ok := mem.Find(internalKey)
	if !ok && imm != nil {
		key, value, ok = imm.Find(internalKey)
		if !ok && fmem != nil {
			key, value, ok = fmem.Find(internalKey)
			if !ok {
				//todo: read sstable
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

	return &proto.ReadResponse{
		Value: value,
		Ts:    key.GetTimeStamp(),
	}, nil
}

func (node *StorageNode) Set(ctx context.Context, request *proto.WriteRequest) (*proto.WriteResponse, error) {
	log.Info("set request",
		zap.String("key", request.Key),
		zap.Uint64("timestamp", request.Ts),
		zap.String("value", request.Value))

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

	internalKey := db.NewInternalKey(request.Key, request.Ts)
	node.mem.Lock(slot)
	mem := node.mem.GetMemTable(slot)
	if mem == nil {
		mem = node.mem.CreateTables(slot)
	}
	mem.Set(internalKey, request.Value)
	node.mem.Unlock(slot)

	return &proto.WriteResponse{}, nil
}

func (node *StorageNode) Remove(ctx context.Context, request *proto.WriteRequest) (*proto.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

// Storage service
// Coordinator to storage nodes
func (node *StorageNode) Compact(ctx context.Context, request *proto.CompactRequest) (*proto.CompactResponse, error) {
	slot := common.SlotID(request.Slot)
	group := node.mem.Rotate(slot)
	imm := group.Imm()

	if request.Leader == node.Info.Id {
		slotCh := node.compactor.PreCompact(slot)
		ch := slotCh.AcquireChan()

		go func() {
			imm.Range(func(key db.InternalKey, value string) bool {
				ch <- &proto.KvData{
					Key:       key.UserKey(),
					Value:     value,
					Timestamp: key.GetTimeStamp(),
				}

				return true
			})

			close(ch)
		}()

		mergeCh := node.compactor.MergeChan(slotCh)

		node.WriteSSTable(mergeCh)
	}

	return &proto.CompactResponse{}, nil
}

// Storage nodes to leader storage node
func (node *StorageNode) PushMemTable(ctx context.Context, request *proto.PushMemTableRequest) (*proto.PushMemTableResponse, error) {
	return nil, nil
}

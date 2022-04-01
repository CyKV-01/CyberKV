package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KeyValue service

func (node *StorageNode) Get(ctx context.Context, request *proto.ReadRequest) (*proto.ReadResponse, error) {
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
	slot := common.CalcSlotID(request.Key)

	node.walMutex.Lock()
	wal, ok := node.wals[slot]
	if !ok {
		wal = NewLogWriter(slot, node.Info.Id)
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
	slot := common.SlotID(request.Slot)
	group := node.mem.Rotate(slot)
	imm := group.Imm()

	if request.Leader == node.Info.Id {
		log.Info("ready to compact memtable as leader",
			zap.Int32("slot", slot))

		log.Info("create slot channel",
			zap.Int32("slot", slot))
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
		if len(mergeCh) == 0 {
			log.Warn("merged channel is empty!")
		}

		err := node.tableMgr.WriteSSTable(ctx, mergeCh, 1)
		if err != nil {
			log.Error("failed to write sstable",
				zap.Int("level", 1),
				zap.Error(err))
			return nil, err
		}

		log.Info("compact done")
	} else {
		log.Info("ready to compact memtable as follower",
			zap.Int32("slot", slot))

		data := make([]*proto.KvData, 0, 100)
		imm.Range(func(key db.InternalKey, value string) bool {
			data = append(data, &proto.KvData{
				Key:       key.UserKey(),
				Value:     value,
				Timestamp: key.GetTimeStamp(),
			})

			return true
		})
		req := proto.PushMemTableRequest{
			Slot: request.Slot,
			Data: data,
		}

		conn, err := grpc.Dial(request.LeaderAddr, grpc.WithInsecure())
		if err != nil {
			log.Error("failed to connect to leader",
				zap.String("leader_addr", request.LeaderAddr),
				zap.Error(err))
			return nil, err
		}
		cli := proto.NewStorageClient(conn)

		_, err = cli.PushMemTable(ctx, &req)
		if err != nil {
			log.Error("failed to push memtable to leader",
				zap.String("leader_addr", request.LeaderAddr),
				zap.Error(err))
			return nil, err
		}
	}

	return &proto.CompactMemTableResponse{}, nil
}

// Storage nodes to leader storage node
func (node *StorageNode) PushMemTable(ctx context.Context, request *proto.PushMemTableRequest) (*proto.PushMemTableResponse, error) {
	slotCh := node.compactor.GetChan(request.Slot)
	for slotCh == nil {
		time.Sleep(time.Second)
		slotCh = node.compactor.GetChan(request.Slot)
	}
	ch := slotCh.AcquireChan()

	for _, data := range request.Data {
		ch <- data
	}

	close(ch)

	return &proto.PushMemTableResponse{}, nil
}

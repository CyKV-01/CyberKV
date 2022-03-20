package storage

import (
	"context"
	"fmt"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (node *StorageNode) Get(ctx context.Context, request *proto.ReadRequest) (*proto.ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func (node *StorageNode) Set(ctx context.Context, request *proto.WriteRequest) (*proto.WriteResponse, error) {
	slot := common.CalcSlotID(request.Key)

	node.walMutex.Lock()
	wal, ok := node.wals[slot]
	if !ok {
		wal = NewLogWritter(slot)
		node.wals[slot] = wal
	}
	node.walMutex.Unlock()

	batch := common.NewBatch()
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

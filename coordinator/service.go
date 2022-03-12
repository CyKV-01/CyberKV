package coordinator

import (
	"context"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
)

func (coord *Coordinator) Get(ctx context.Context, request *proto.ReadRequest) (response *proto.ReadResponse, err error) {
	slot := common.CalcSlotID(request.Key)

	nodes := coord.computeCluster.GetNodesBySlot(slot)
	if len(nodes) == 0 {
		return &proto.ReadResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "No node to serve",
			},
		}, nil
	}

	node := nodes[0].(*ComputeNode)

	return node.Get(ctx, request)
}

func (coord *Coordinator) Set(ctx context.Context, request *proto.WriteRequest) (response *proto.WriteResponse, err error) {
	slot := common.CalcSlotID(request.Key)

	nodes := coord.computeCluster.GetNodesBySlot(slot)
	if len(nodes) == 0 {
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "No node to serve",
			},
		}, nil
	}

	node := nodes[0].(*ComputeNode)

	return node.Set(ctx, request)
}

func (coord *Coordinator) Remove(ctx context.Context, request *proto.WriteRequest) (response *proto.WriteResponse, err error) {
	slot := common.CalcSlotID(request.Key)

	nodes := coord.computeCluster.GetNodesBySlot(slot)
	if len(nodes) == 0 {
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "No node to serve",
			},
		}, nil
	}

	node := nodes[0].(*ComputeNode)

	return node.Remove(ctx, request)
}

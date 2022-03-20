package coordinator

import (
	"context"
	"time"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
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

	node := nodes[0]

	resp, err := node.Get(ctx, request)
	if err != nil {
		log.Errorf("failed to get for key=%s, err=%v", request.Key, err)
		return nil, err
	}

	return resp, nil
}

func (coord *Coordinator) Set(ctx context.Context, request *proto.WriteRequest) (response *proto.WriteResponse, err error) {
	slot := common.CalcSlotID(request.Key)
	nodes := coord.computeCluster.GetNodesBySlot(slot)
	if nodes == nil {
		log.Warn("hit a new slot, schedule it...")
		slots := []common.SlotID{slot}
		coord.computeCluster.AssignSlots(slots)
		time.Sleep(50 * time.Millisecond)
		nodes = coord.computeCluster.GetNodesBySlot(slot)
	}
	if len(nodes) == 0 {
		log.Warn("no node to serve")
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "No node to serve",
			},
		}, nil
	}

	node := nodes[0]

	storageNodes := coord.storageCluster.GetNodesBySlot(slot)
	if storageNodes == nil {
		log.Warn("hit a new slot, schedule it...")
		slots := []common.SlotID{slot}
		coord.storageCluster.AssignSlots(slots)
		time.Sleep(50 * time.Millisecond)
		storageNodes = coord.storageCluster.GetNodesBySlot(slot)
	}
	if len(storageNodes) == 0 {
		log.Warn("no storage node to serve")
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "No storage node to serve",
			},
		}, nil
	}
	for _, node := range storageNodes {
		request.Info = append(request.Info, &node.info.NodeInfo)
	}

	resp, err := node.Set(ctx, request)
	if err != nil {
		log.Errorf("failed to set for key=%s value=%s, err=%v",
			request.Key, request.Value, err)
		return nil, err
	}

	return resp, nil
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

	node := nodes[0]

	return node.Remove(ctx, request)
}

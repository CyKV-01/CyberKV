package coordinator

import (
	"context"
	"sync/atomic"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
)

func (coord *Coordinator) Get(ctx context.Context, request *proto.ReadRequest) (response *proto.ReadResponse, err error) {
	if request.Ts <= 0 {
		request.Ts = coord.CurrentTs()
	}

	slot := common.CalcSlotID(request.Key)
	nodes := coord.computeCluster.GetNodesBySlot(slot)
	if len(nodes) < coord.computeCluster.readQuorum {
		log.Warn("no enough compute node to serve",
			zap.Int("nodes_num", len(nodes)),
			zap.Int("read_quorum", coord.computeCluster.readQuorum))
		return &proto.ReadResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "no enough compute node to serve",
			},
		}, nil
	}

	node := nodes[0]

	storageNodes := coord.storageCluster.GetNodesBySlot(slot)
	if len(storageNodes) == 0 {
		log.Warn("no enough storage node to serve",
			zap.Int("nodes_num", len(nodes)),
			zap.Int("read_quorum", coord.computeCluster.readQuorum))
		return &proto.ReadResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "no enough storage node to serve",
			},
		}, nil
	}
	for _, node := range storageNodes {
		request.Info = append(request.Info, &node.info.NodeInfo)
	}

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
	if len(nodes) < coord.computeCluster.writeQuorum {
		log.Warn("no enough compute node to serve",
			zap.Int("nodesNum", len(nodes)),
			zap.Int("writeQuorum", coord.computeCluster.writeQuorum))
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "no enough compute node to serve",
			},
		}, nil
	}

	node := nodes[0]

	storageNodes := coord.storageCluster.GetNodesBySlot(slot)
	if len(storageNodes) == 0 {
		log.Warn("no enough storage node to serve",
			zap.Int("nodesNum", len(nodes)),
			zap.Int("writeQuorum", coord.computeCluster.writeQuorum))
		return &proto.WriteResponse{
			Status: &proto.Status{
				ErrCode:    proto.ErrorCode_RetryLater,
				ErrMessage: "no enough storage node to serve",
			},
		}, nil
	}
	for _, node := range storageNodes {
		request.Info = append(request.Info, &node.info.NodeInfo)
	}

	request.Ts = coord.GenTs()

	resp, err := node.Set(ctx, request)
	if err != nil {
		log.Errorf("failed to set for key=%s value=%s, err=%v",
			request.Key, request.Value, err)
		return nil, err
	}

	if common.IsOk(resp.Status) {
		old := coord.AddSlotMemSize(slot, uint64(len(request.Key)+len(request.Value)))
		if old >= common.WalCompactThreshold &&
			atomic.CompareAndSwapUint64(&coord.slotMemSizeTable[slot], old, 0) {
			coord.CompactMemTable(slot)
		}
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

// Coordinator service
func (coord *Coordinator) AllocateSSTableID(ctx context.Context, request *proto.AllocateSSTableRequest) (*proto.AllocateSSTableResponse, error) {
	id, err := coord.sstableIdAllocator.Next()
	if err != nil {
		return nil, err
	}

	return &proto.AllocateSSTableResponse{
		Id: id,
	}, nil
}

func (coord *Coordinator) ReportStats(ctx context.Context, request *proto.ReportStatsRequest) (*proto.ReportStatsResponse, error) {
	_, ok := coord.storageCluster.GetNode(request.Id)
	if !ok { // ignore
		return &proto.ReportStatsResponse{}, nil
	}

	maxSize := uint64(0)
	for slot, size := range request.MemTableSize {
		if size > common.WalCompactThreshold {
			log.Info("schedule memtable compaction",
				zap.Int32("slot", slot),
				zap.Uint64("mem_table_size", size))
			nodes := coord.storageCluster.GetNodesBySlot(slot)
			if len(nodes) != coord.storageCluster.replicaNum {
				log.Warn("the number of nodes will compact is not equal to the replica number",
					zap.Int("nodesNum", len(nodes)),
					zap.Int("replicaNum", coord.storageCluster.replicaNum))
			}
			go coord.compactor.CompactMemTable(slot, nodes)
		}
		if size > maxSize {
			maxSize = size
		}
	}

	log.Info("receive report",
		zap.String("node_id", request.Id),
		zap.Uint64("max_mem_table_size", maxSize))

	return &proto.ReportStatsResponse{}, nil
}

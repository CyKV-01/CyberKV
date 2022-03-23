package common

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Node interface {
	Register(name string)
	Start()
	Recovery()
}

type BaseNode struct {
	Meta *etcdcli.Client
	Info *proto.NodeInfo
}

func NewBaseNode(addr string, etcd *etcdcli.Client) *BaseNode {
	info := &proto.NodeInfo{
		Id:   uuid.NewString(),
		Addr: addr,
	}

	return &BaseNode{
		Info: info,
		Meta: etcd,
	}
}

func (node *BaseNode) Register(name string) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)

	resp, err := node.Meta.Grant(ctx, DefaultTTL)
	if err != nil {
		log.Errorf("failed to grant a lease, err=%v", err)
		panic(err)
	}

	keepAliveCtx := context.Background()
	keepAliveStream, err := node.Meta.KeepAlive(keepAliveCtx, resp.ID)
	if err != nil {
		log.Errorf("failed to keep alive a lease, err=%v", err)
		panic(err)
	}
	go func() {
		for {
			select {
			case resp := <-keepAliveStream:
				log.Debug("etcd keep alive response",
					zap.Int64("id", int64(resp.ID)),
					zap.Int64("ttl", resp.TTL))
			case <-keepAliveCtx.Done():
				log.Infof("etcd keep alive canceled",
					zap.String("nodeId", node.Info.Id))
				return
			}
		}
	}()

	infoBytes, err := json.Marshal(node.Info)
	if err != nil {
		log.Errorf("failed to marshal info, err=%v", err)
		panic(err)
	}

	log.Info("register",
		zap.String("node", name),
		zap.String("id", node.Info.Id),
		zap.String("addr", node.Info.Addr))
	_, err = node.Meta.Put(ctx, path.Join(ServicePrefix, name, node.Info.Id), string(infoBytes),
		etcdcli.WithLease(resp.ID))
	if err != nil {
		log.Errorf("failed to register, err=%v", err)
		panic(err)
	}
}

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
	Register()
	Start()
}

type BaseNode struct {
	Etcd *etcdcli.Client
	Info *proto.NodeInfo
}

func NewBaseNode(addr string, etcd *etcdcli.Client) *BaseNode {
	info := &proto.NodeInfo{
		Id:   uuid.NewString(),
		Addr: addr,
	}

	return &BaseNode{
		Info: info,
		Etcd: etcd,
	}
}

func (node *BaseNode) Register() {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)

	resp, err := node.Etcd.Grant(ctx, DefaultTTL)
	if err != nil {
		log.Errorf("failed to grant a lease, err=%v", err)
		panic(err)
	}

	_, err = node.Etcd.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		log.Errorf("failed to keep alive a lease, err=%v", err)
		panic(err)
	}

	infoBytes, err := json.Marshal(node.Info)
	if err != nil {
		log.Errorf("failed to marshal info, err=%v", err)
		panic(err)
	}

	log.Info("register coordinator",
		zap.String("id", node.Info.Id),
		zap.String("addr", node.Info.Addr))
	_, err = node.Etcd.Put(ctx, path.Join(ServicePrefix, "coordinator", node.Info.Id), string(infoBytes),
		etcdcli.WithLease(resp.ID))
	if err != nil {
		log.Errorf("failed to register, err=%v", err)
		panic(err)
	}
}

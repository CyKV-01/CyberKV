package common

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Component interface {
	Register(name string)
	Start()
	Recover()
}

type BaseComponent struct {
	Meta *etcdcli.Client
	Info *proto.NodeInfo
}

func NewBaseComponent(addr string, etcd *etcdcli.Client) *BaseComponent {
	id, err := GLobalSonyflake.NextID()
	if err != nil {
		panic(err)
	}

	info := &proto.NodeInfo{
		Id:   id,
		Addr: addr,
	}

	return &BaseComponent{
		Info: info,
		Meta: etcd,
	}
}

func (node *BaseComponent) Recover() {
	nodeInfoBytes, err := ioutil.ReadFile("component.json")
	hasComponentInfo := false
	if err != nil {
		log.Warn("failed to read component infomation from disk",
			zap.Error(err))
	} else {
		err = json.Unmarshal(nodeInfoBytes, node.Info)
		if err != nil {
			log.Warn("failed to deserialize component infomation",
				zap.Error(err))
		} else {
			hasComponentInfo = true
		}
	}

	if !hasComponentInfo {
		file, err := os.Create("component.json")
		if err != nil {
			panic(err)
		}

		infoBytes, err := json.Marshal(node.Info)
		if err != nil {
			panic(err)
		}
		_, err = file.Write(infoBytes)
		if err != nil {
			panic(err)
		}

		err = file.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (node *BaseComponent) Register(name string) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 2*time.Second)

	log.Info("register...",
		zap.String("name", name))

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
			case <-keepAliveStream:
				// log.Debug("etcd keep alive response",
				// 	zap.Int64("id", int64(resp.ID)),
				// 	zap.Int64("ttl", resp.TTL))
			case <-keepAliveCtx.Done():
				log.Infof("etcd keep alive canceled",
					zap.Uint64("nodeId", node.Info.Id))
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
		zap.Uint64("id", node.Info.Id),
		zap.String("addr", node.Info.Addr))
	_, err = node.Meta.Put(ctx, path.Join(ServicePrefix, name, strconv.FormatUint(node.Info.Id, 10)),
		string(infoBytes),
		etcdcli.WithLease(resp.ID))
	if err != nil {
		log.Errorf("failed to register, err=%v", err)
		panic(err)
	}
}

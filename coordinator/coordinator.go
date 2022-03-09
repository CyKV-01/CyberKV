package coordinator

import (
	"context"

	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Coordinator struct {
	proto.UnimplementedKeyValueServer

	etcdClient *etcdcli.Client

	computeCluster *Cluster
	storageCluster *Cluster
}

func (coord *Coordinator) Get(context.Context, *proto.ReadRequest) (*proto.ReadResponse, error) {
	ComputeCluster
}
func (coord *Coordinator) Set(context.Context, *proto.WriteRequest) (*proto.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (coord *Coordinator) Remove(context.Context, *proto.WriteRequest) (*proto.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

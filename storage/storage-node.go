package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/yah01/CyberKV/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const SSTableRootDir = "data"

type StorageNode struct {
	proto.UnimplementedKeyValueServer

	minioClient *minio.Client
}

func NewStorageNode(minioClient *minio.Client) *StorageNode {
	return &StorageNode{
		minioClient: minioClient,
	}
}

func (node *StorageNode) Get(context.Context, *proto.ReadRequest) (*proto.ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (node *StorageNode) Set(context.Context, *proto.WriteRequest) (*proto.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (node *StorageNode) Remove(context.Context, *proto.WriteRequest) (*proto.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

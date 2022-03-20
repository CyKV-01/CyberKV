package main

import (
	"github.com/minio/minio-go/v7"
	"github.com/yah01/CyberKV/storage"
	etcdcli "go.etcd.io/etcd/client/v3"
)

func main() {
	etcdClient, err := etcdcli.New(etcdcli.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}

	minioClient, err := minio.New("127.0.0.1:9000", &minio.Options{})
	if err != nil {
		panic(err)
	}

	storageNode := storage.NewStorageNode("[::]:4568", etcdClient, minioClient)
	storageNode.Register()
	storageNode.Start()
}

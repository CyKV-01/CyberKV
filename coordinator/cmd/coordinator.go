package main

import (
	"github.com/yah01/CyberKV/coordinator"
	etcdcli "go.etcd.io/etcd/client/v3"
)

func main() {
	etcdClient, err := etcdcli.New(etcdcli.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}

	coord := coordinator.NewCoordinator(etcdClient, "[::]:4567")

	coord.Start()
}

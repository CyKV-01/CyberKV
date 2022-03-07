package coordinator

import etcdcli "go.etcd.io/etcd/client/v3"

type Coordinator struct {
	etcdClient *etcdcli.Client
	
	computeCluster *Cluster
	storageCluster *Cluster
}

type 

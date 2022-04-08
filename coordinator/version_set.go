package coordinator

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/golang/protobuf/proto"
	"github.com/yah01/CyberKV/common"
	. "github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	etcdcli "go.etcd.io/etcd/client/v3"
)

type VersionSet struct {
	guard   sync.RWMutex
	meta    *etcdcli.Client
	current *proto.Version
	last    UniqueID
}

func NewVersionSet(meta *etcdcli.Client) (*VersionSet, error) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	resp, err := meta.Get(ctx, common.VersionSetKey)
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		_, err = meta.Delete(ctx, common.VersionPrefix, etcdcli.WithPrefix())
		if err != nil {
			return nil, err
		}

		id, err := GLobalSonyflake.NextID()
		if err != nil {
			return nil, err
		}
		version := &proto.Version{
			VersionId: id,
			Sstables:  nil,
		}

		log.Debug("save version...")
		err = saveVersion(meta, version)
		if err != nil {
			return nil, err
		}

		versionSet := &proto.VersionSet{
			Current: version.VersionId,
		}

		log.Debug("save version set...")
		saveVersionSet(meta, versionSet)
		if err != nil {
			return nil, err
		}

		return &VersionSet{
			guard:   sync.RWMutex{},
			meta:    meta,
			current: version,
		}, nil
	}

	var (
		versionSet proto.VersionSet
		version    proto.Version
	)
	err = pb.Unmarshal(resp.Kvs[0].Value, &versionSet)
	if err != nil {
		return nil, err
	}

	resp, err = meta.Get(ctx, fmt.Sprintf("%s/%v", VersionPrefix, versionSet.Current))
	if err != nil {
		return nil, err
	}

	err = pb.Unmarshal(resp.Kvs[0].Value, &version)
	if err != nil {
		return nil, err
	}

	return &VersionSet{
		guard:   sync.RWMutex{},
		meta:    meta,
		current: &version,
		last:    versionSet.Last,
	}, nil
}

func (set *VersionSet) Current() *proto.Version {
	set.guard.RLock()
	defer set.guard.RUnlock()

	return set.current
}

func (set *VersionSet) Edit(addition []*proto.SSTableLevel, deletion []*proto.SSTableLevel) (*proto.Version, error) {
	version := pb.Clone(set.current).(*proto.Version)

	sstables := make(map[int32]Set[string])
	for level, tables := range version.Sstables {
		levelTable := make(Set[string])
		levelTable.Insert(tables.Sstables...)
		sstables[level] = levelTable
	}

	for _, del := range deletion {
		levelTable, ok := sstables[del.Level]
		if ok {
			levelTable.Delete(del.Sstables...)
		}
	}

	for _, add := range addition {
		levelTable, ok := sstables[add.Level]
		if ok {
			levelTable.Insert(add.Sstables...)
		}
	}

	id, err := GLobalSonyflake.NextID()
	if err != nil {
		return nil, err
	}
	version.VersionId = id
	for level, tables := range sstables {
		version.Sstables[level] = &proto.SSTableLevel{
			Level:    level,
			Sstables: tables.Collect(),
		}
	}

	err = saveVersion(set.meta, version)
	if err != nil {
		return nil, err
	}
	err = saveVersionSet(set.meta, &proto.VersionSet{
		Current: version.VersionId,
		Last:    set.current.VersionId,
	})
	if err != nil {
		return nil, err
	}

	set.guard.Lock()
	set.last = set.current.VersionId
	set.current = version
	set.guard.Unlock()

	return pb.Clone(version).(*proto.Version), nil
}

func saveVersion(meta *etcdcli.Client, version *proto.Version) error {
	ctx := context.Background()

	versionBytes, err := pb.Marshal(version)
	if err != nil {
		return err
	}

	_, err = meta.Put(ctx, fmt.Sprintf("%s/%v", VersionPrefix, version.VersionId), string(versionBytes))
	if err != nil {
		return err
	}

	return nil
}

func saveVersionSet(meta *etcdcli.Client, versionSet *proto.VersionSet) error {
	ctx := context.Background()

	versionSetBytes, err := pb.Marshal(versionSet)
	if err != nil {
		return err
	}

	_, err = meta.Put(ctx, VersionSetKey, string(versionSetBytes))
	if err != nil {
		return err
	}

	return nil
}

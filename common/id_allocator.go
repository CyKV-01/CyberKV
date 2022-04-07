package common

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/sony/sonyflake"
	etcdcli "go.etcd.io/etcd/client/v3"
)

type IdAllocator interface {
	NextID() (uint64, error)
}

type MetaIdAllocator struct {
	meta      *etcdcli.Client
	key       string
	currentId uint64
	endId     uint64
	batchSize uint64
}

func InitMetaAllocator(meta *etcdcli.Client, key string) {
	var err error
	GlobalMetaIdAllocator, err = NewMetaIdAllocator(context.Background(), meta, key, 100)
	if err != nil {
		panic(err)
	}
}

func NewMetaIdAllocator(ctx context.Context, meta *etcdcli.Client, key string, batchSize uint64) (*MetaIdAllocator, error) {
	resp, err := meta.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	id := uint64(1)
	if resp.Count > 0 {
		id, err = strconv.ParseUint(string(resp.Kvs[0].Value), 10, 64)
		if err != nil {
			return nil, err
		}
	}

	endId := id + batchSize

	_, err = meta.Put(ctx, key, fmt.Sprint(endId))
	if err != nil {
		return nil, err
	}

	return &MetaIdAllocator{
		meta:      meta,
		key:       key,
		currentId: id,
		endId:     endId,
		batchSize: batchSize,
	}, nil
}

func (allocator *MetaIdAllocator) NextID() (uint64, error) {
	if allocator.currentId == allocator.endId {
		var (
			ctx   = context.Background()
			id    uint64
			endId uint64
		)

		for {
			resp, err := allocator.meta.Get(ctx, allocator.key)
			if err != nil {
				return 0, err
			}

			id, err := strconv.ParseUint(string(resp.Kvs[0].Value), 10, 64)
			if err != nil {
				return 0, err
			}
			endId = id + allocator.batchSize

			txn := allocator.meta.Txn(ctx)
			txn.If(etcdcli.Compare(etcdcli.Value(allocator.key), "=", resp.Kvs[0].Value)).
				Then(
					etcdcli.OpPut(allocator.key, fmt.Sprint(endId)),
				)

			_, err = txn.Commit()
			if err == nil {
				break
			}
		}

		allocator.currentId = id
		allocator.endId = endId
	}

	return atomic.AddUint64(&allocator.currentId, 1) - 1, nil
}

var (
	Sonyflake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			return uint16(time.Now().Nanosecond() % 10007), nil
		},
	})

	GlobalMetaIdAllocator *MetaIdAllocator
)

package common

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	etcdcli "go.etcd.io/etcd/client/v3"
)

type IdAllocator struct {
	meta      *etcdcli.Client
	key       string
	currentId uint64
	endId     uint64
	batchSize uint64
}

func NewIdAllocator(ctx context.Context, meta *etcdcli.Client, key string, batchSize uint64) (*IdAllocator, error) {
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

	return &IdAllocator{
		meta:      meta,
		key:       key,
		currentId: id,
		endId:     endId,
		batchSize: batchSize,
	}, nil
}

func (allocator *IdAllocator) Next() (uint64, error) {
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

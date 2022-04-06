package storage

import (
	"sync"

	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/proto"
)

const (
	DefaultKvDataBufferSize = common.WalCompactThreshold / 10
)

var (
	kvDataBufferPool = sync.Pool{
		New: func() any {
			return make([]*proto.KvData, 0, DefaultKvDataBufferSize)
		},
	}
)

func GetKvDataBuffer() []*proto.KvData {
	return kvDataBufferPool.Get().([]*proto.KvData)[:0]
}

func PutBackToPool(buffer []*proto.KvData) {
	kvDataBufferPool.Put(buffer)
}

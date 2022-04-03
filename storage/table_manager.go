package storage

import (
	"bytes"
	"context"
	"io/ioutil"
	"sort"

	pb "github.com/golang/protobuf/proto"
	"github.com/minio/minio-go/v7"
	"github.com/yah01/CyberKV/common"
	"github.com/yah01/CyberKV/common/db"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
)

type TableManager struct {
	bucketName string
	store      *minio.Client
	coord      proto.CoordinatorClient
}

func NewTableManager(bucketName string, store *minio.Client, coord proto.CoordinatorClient) *TableManager {
	return &TableManager{
		bucketName: bucketName,
		store:      store,
		coord:      coord,
	}
}

// Todo: get sstable and index concurrently
func (mgr *TableManager) Get(ctx context.Context, key db.InternalKey) (value string, ok bool, err error) {
	for level := 1; level <= common.MaxLevel; level++ {
		objects := mgr.store.ListObjects(ctx, mgr.bucketName, minio.ListObjectsOptions{
			WithMetadata: true,
			Prefix:       common.SSTableDataLevelPrefix(level),
		})

		for objectInfo := range objects {
			var (
				sstableObject *minio.Object
				sstableData   []byte

				indexObject *minio.Object
				indexData   []byte
			)
			if !mgr.CheckBloomFilter(key.UserKey(), objectInfo.UserMetadata) {
				continue
			}

			if key.UserKey() < objectInfo.UserMetadata["first"] ||
				key.UserKey() > objectInfo.UserMetadata["last"] {
				continue
			}

			sstableObject, err = mgr.store.GetObject(ctx, mgr.bucketName, objectInfo.Key, minio.GetObjectOptions{})
			if err != nil {
				return
			}

			indexObject, err = mgr.store.GetObject(ctx, mgr.bucketName, common.GetSSTableIndexPath(objectInfo.Key), minio.GetObjectOptions{})
			if err != nil {
				return
			}

			sstableData, err = ioutil.ReadAll(sstableObject)
			if err != nil {
				return
			}

			indexData, err = ioutil.ReadAll(indexObject)
			if err != nil {
				return
			}

			var index proto.Index
			err = pb.UnmarshalMerge(indexData, &index)
			if err != nil {
				return
			}

			sstable := db.NewSSTable(sstableData, level, objectInfo.Key, objectInfo.UserMetadata["first"], objectInfo.UserMetadata["last"], &index)

			value, ok = sstable.Get(key.UserKey(), key.GetTimeStamp())
			if ok {
				return
			}
		}
	}

	return
}

// Write SSTable to object storage
// This will remove all overlapping SSTables
func (mgr *TableManager) WriteSSTable(ctx context.Context, dataCh chan *proto.KvData, level int) ([]*db.SSTable, []*db.SSTable, error) {
	log.Debug("create sstable from merged data channel")
	sstable := db.NewSSTableFromDataCh(dataCh)
	if sstable.Empty() {
		log.Warn("data is empty, don't write")
		return nil, nil, nil
	}

	overlapTables := mgr.GetOverlapTables(ctx, sstable.FirstKey, sstable.LastKey, level)

	mergedDataCh := mgr.MergeSSTables(sstable, overlapTables)

	log.Info("generate new sstable...")
	newSSTable := db.NewSSTableFromDataCh(mergedDataCh)

	err := mgr.writeSSTable(ctx, newSSTable, level)
	if err != nil {
		return nil, nil, err
	}

	createdSSTables := make([]*db.SSTable, 0, 1)
	createdSSTables = append(createdSSTables, newSSTable)

	return createdSSTables, overlapTables, nil
}

// Write directly the SSTable into object storage
func (mgr *TableManager) writeSSTable(ctx context.Context, sstable *db.SSTable, level int) error {
	log.Info("write sstable into object storage...",
		zap.Int("size", sstable.Len()),
		zap.String("firstKey", sstable.FirstKey),
		zap.String("lastKey", sstable.LastKey))

	idResp, err := mgr.coord.AllocateSSTableID(ctx, &proto.AllocateSSTableRequest{})
	if err != nil {
		return err
	}

	sstable.Level = level
	sstable.Path = common.SSTableDataPath(level, idResp.Id)
	_, err = mgr.store.PutObject(ctx, mgr.bucketName, sstable.Path,
		bytes.NewReader(sstable.GetData()), int64(sstable.Len()), minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"first": sstable.FirstKey,
				"last":  sstable.LastKey,
			},
		})
	if err != nil {
		return err
	}

	indexBytes, err := pb.Marshal(sstable.GetIndex())
	if err != nil {
		return err
	}

	_, err = mgr.store.PutObject(ctx, mgr.bucketName, common.GetSSTableIndexPath(sstable.Path),
		bytes.NewReader(indexBytes), int64(len(indexBytes)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (mgr *TableManager) MergeSSTables(src *db.SSTable, overlaps []*db.SSTable) chan *proto.KvData {
	log.Info("merge sstables...")
	mergeCh := make(chan *proto.KvData, 32)

	go func() {
		srcIter := db.NewSSTableIter(src)
		for _, overlap := range overlaps {
			overlapIter := db.NewSSTableIter(overlap)

			for ; overlapIter.Valid(); overlapIter.Next() {
				overlapKv := overlapIter.Value()
				overlapInternalKey := db.NewInternalKey(overlapKv.Key, overlapKv.Timestamp)

				for ; srcIter.Valid(); srcIter.Next() {
					srcKv := srcIter.Value()
					srcInternalKey := db.NewInternalKey(srcKv.Key, srcKv.Timestamp)

					if srcInternalKey.Compare(overlapInternalKey) < 0 {
						log.Info("insert data from src",
							zap.String("key", srcKv.Key))
						mergeCh <- &proto.KvData{
							Key:       srcKv.Key,
							Value:     srcKv.Value,
							Timestamp: srcKv.Timestamp,
						}
					} else {
						break
					}
				}

				mergeCh <- &proto.KvData{
					Key:       overlapKv.Key,
					Value:     overlapKv.Value,
					Timestamp: overlapKv.Timestamp,
				}
			}
		}

		// There may be data in src left
		for ; srcIter.Valid(); srcIter.Next() {
			srcKv := srcIter.Value()

			mergeCh <- &proto.KvData{
				Key:       srcKv.Key,
				Value:     srcKv.Value,
				Timestamp: srcKv.Timestamp,
			}
		}

		close(mergeCh)
	}()

	return mergeCh
}

func (mgr *TableManager) GetOverlapTables(ctx context.Context, firstKey, lastKey string, level int) []*db.SSTable {
	log.Info("fetch overlapping sstables",
		zap.String("firstKey", firstKey),
		zap.String("lastKey", lastKey))

	objectCh := mgr.store.ListObjects(ctx, mgr.bucketName, minio.ListObjectsOptions{
		WithMetadata: true,
		Prefix:       common.SSTableDataLevelPrefix(level),
	})

	sstables := make([]*db.SSTable, 0, len(objectCh))
	for object := range objectCh {
		if object.UserMetadata["first"] > lastKey || object.UserMetadata["last"] < firstKey {
			continue
		}

		sstables = append(sstables,
			db.NewSSTable(nil, level, object.Key, object.UserMetadata["first"], object.UserMetadata["last"], nil))

		log.Debug("overlap sstable for given key range",
			zap.String("firstKey", firstKey),
			zap.String("lastKey", lastKey),
			zap.String("sstable", object.Key))
	}

	sort.Slice(sstables, func(i, j int) bool {
		return sstables[i].FirstKey < sstables[j].FirstKey
	})

	return sstables
}

func (mgr *TableManager) CheckBloomFilter(key string, filter map[string]string) bool {
	return true
}

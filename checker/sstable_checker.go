package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yah01/CyberKV/common/db"
)

func main() {
	minioEndpoint := "127.0.0.1:9000"
	store, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	sstableObject, err := store.GetObject(ctx, "cyberkv", "sstable/0/7601.cdb", minio.GetObjectOptions{})
	if err != nil {
		return
	}

	sstableData, err := ioutil.ReadAll(sstableObject)
	if err != nil {
		return
	}

	sstable := db.NewSSTable(sstableData, 0, "sstable/1/4101.cdb", "", "", nil)

	iter := db.NewSSTableIter(sstable)
	for ; iter.Valid(); iter.Next() {
		fmt.Printf("key=%s value=%s ts=%v\n",
			iter.Value().Key, iter.Value().Value, iter.Value().Timestamp)
	}
}

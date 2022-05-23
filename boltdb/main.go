package main

import (
	"context"
	"net"

	"github.com/boltdb/bolt"
	"github.com/yah01/CyberKV/proto"
	"google.golang.org/grpc"
)

type BoltDB struct {
	proto.UnimplementedKeyValueServer
	inner *bolt.DB
}

const BucketName = "cyberkv-boltdb"

func NewBoltDB() *BoltDB {
	db, err := bolt.Open("bench-test", 0600, nil)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(BucketName))
		if err != nil {
			panic(err)
		}

		return nil
	})

	return &BoltDB{
		inner: db,
	}
}

func (db *BoltDB) AssignSlot(ctx context.Context, req *proto.AssignSlotRequest) (*proto.AssignSlotResponse, error) {
	return nil, nil
}
func (db *BoltDB) Get(ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {
	var value string
	db.inner.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		value = string(bucket.Get([]byte(req.Key)))

		return nil
	})

	return &proto.ReadResponse{
		Value: value,
	}, nil
}
func (db *BoltDB) Set(ctx context.Context, req *proto.WriteRequest) (*proto.WriteResponse, error) {
	err := db.inner.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		return bucket.Put([]byte(req.Key), []byte(req.Value))
	})

	if err != nil {
		return nil, err
	}

	return &proto.WriteResponse{}, nil
}
func (db *BoltDB) Remove(ctx context.Context, req *proto.WriteRequest) (*proto.WriteResponse, error) {
	return &proto.WriteResponse{}, nil
}

func main() {
	db := NewBoltDB()

	server := grpc.NewServer()
	proto.RegisterKeyValueServer(server, db)

	listener, err := net.Listen("tcp", ":9595")
	if err != nil {
		panic(err)
	}

	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

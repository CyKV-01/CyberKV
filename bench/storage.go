package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/yah01/CyberKV/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4568", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cli := proto.NewKeyValueClient(conn)

	value := "0"
	for i := 0; i < 4000; i++ {
		value += "0"
	}
	req := proto.WriteRequest{
		Key:   "hello",
		Value: "world",
	}

	ctx := context.Background()
	cnt := 10000
	concurrency := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}
	begin := time.Now()
	for c := 0; c < concurrency; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < cnt; i++ {
				cli.Set(ctx, &req)
			}
		}()
	}
	wg.Wait()
	end := time.Now()

	duration := end.Sub(begin)
	fmt.Printf("TPS: %v", float64(cnt*concurrency)/duration.Seconds())
}

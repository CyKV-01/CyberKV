package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("[::]:5790", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cli := proto.NewKeyValueClient(conn)

	cnt := 20000
	concurrency := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}
	begin := time.Now()
	for c := 0; c < concurrency; c++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			begin := time.Now()
			for i := 0; i < cnt; i++ {
				if i%100 == 0 {
					end := time.Now()
					msg := fmt.Sprintf("writing data... progress: %d/%d", i, cnt)
					log.Info(msg,
						zap.Int("goroutine", idx),
						zap.Float64("TPS", 100/end.Sub(begin).Seconds()))
					begin = end
				}

				ctx := context.Background()
				ctx, _ = context.WithTimeout(ctx, 5*time.Second)

				req := &proto.WriteRequest{
					Key:   fmt.Sprintf("key-%d", i%10000),
					Value: fmt.Sprintf("value-%d", i),
				}

				_, err := cli.Set(ctx, req)
				if err != nil {
					log.Warn("failed to set",
						zap.Int("goroutine", idx),
						zap.String("key", req.Key),
						zap.String("value", req.Value),
						zap.Error(err))
				}
			}
		}(c)
	}
	wg.Wait()
	end := time.Now()

	duration := end.Sub(begin)
	fmt.Printf("TPS: %v", float64(cnt*concurrency)/duration.Seconds())
}

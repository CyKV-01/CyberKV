package main

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
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

	keyCount := 10000
	cnt := 10000
	outputGap := 1000

	keys := make([]string, keyCount)
	values := make([]string, cnt)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	for i := range values {
		values[i] = "value-" + strconv.Itoa(i)
	}

	concurrency := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}
	begin := time.Now()
	for c := 0; c < concurrency; c++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			var (
				err error
				req = &proto.WriteRequest{}
			)

			begin := time.Now()
			for i := 0; i < cnt; i++ {
				if i > 0 && i%outputGap == 0 {
					end := time.Now()
					msg := fmt.Sprintf("writing data... progress: %d/%d", i, cnt)
					log.Info(msg,
						zap.Int("goroutine", idx),
						zap.Float64("TPS", float64(outputGap)/end.Sub(begin).Seconds()),
						zap.Float64("legency", float64(end.Sub(begin).Milliseconds())/float64(outputGap)))
					begin = end
				}

				ctx := context.Background()
				ctx, _ = context.WithTimeout(ctx, 5*time.Second)

				req.Key = keys[i%keyCount]
				req.Value = values[i]

				_, err = cli.Set(ctx, req)
				if err != nil {
					log.Warn("failed to set",
						zap.Int("goroutine", idx),
						zap.String("key", req.Key),
						zap.String("value", req.Value),
						zap.Error(err))
				}
			}
			end := time.Now()
			msg := fmt.Sprintf("writing data... progress: %d/%d", cnt, cnt)
			log.Info(msg,
				zap.Int("goroutine", idx),
				zap.Float64("TPS", float64(outputGap)/end.Sub(begin).Seconds()))
		}(c)
	}
	wg.Wait()
	end := time.Now()
	duration := end.Sub(begin)
	tps := float64(cnt*concurrency) / duration.Seconds()

	conn, err = grpc.Dial("[::]:5678", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli = proto.NewKeyValueClient(conn)
	begin = time.Now()
	for c := 0; c < concurrency; c++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			var (
				err error
				req = &proto.ReadRequest{}
			)

			begin := time.Now()
			for i := 0; i < cnt; i++ {
				if i > 0 && i%outputGap == 0 {
					end := time.Now()
					msg := fmt.Sprintf("reading data... progress: %d/%d", i, cnt)
					log.Info(msg,
						zap.Int("goroutine", idx),
						zap.Float64("QPS", float64(outputGap)/end.Sub(begin).Seconds()),
						zap.Float64("legency", float64(end.Sub(begin).Milliseconds())/float64(outputGap)))
					begin = end
				}

				ctx := context.Background()
				ctx, _ = context.WithTimeout(ctx, 5*time.Second)

				req = &proto.ReadRequest{
					Key: fmt.Sprintf("key-%d", i%10000),
				}

				_, err = cli.Get(ctx, req)
				if err != nil {
					log.Warn("failed to get",
						zap.Int("goroutine", idx),
						zap.String("key", req.Key),
						zap.Error(err))
				}
			}
			end := time.Now()
			msg := fmt.Sprintf("reading data... progress: %d/%d", cnt, cnt)
			log.Info(msg,
				zap.Int("goroutine", idx),
				zap.Float64("QPS", float64(outputGap)/end.Sub(begin).Seconds()))
		}(c)
	}
	wg.Wait()
	end = time.Now()
	duration = end.Sub(begin)
	qps := float64(cnt*concurrency) / duration.Seconds()

	fmt.Printf("TPS: %v, QPS: %v", tps, qps)
}

package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/urfave/cli/v2"
	"github.com/yah01/CyberKV/storage"
	etcdcli "go.etcd.io/etcd/client/v3"
)

func initPprofMonitor(port int) {
	addr := ":" + strconv.Itoa(port)

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
	}()
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "etcd_endpoints",
				Aliases: []string{"e", "etcd"},
				EnvVars: []string{"ETCD_ENDPOINTS"},
				Value:   cli.NewStringSlice("127.0.0.1:2379"),
			},
			&cli.StringFlag{
				Name:    "minio_endpoint",
				Aliases: []string{"m", "minio"},
				EnvVars: []string{"MINIO_ENDPOINT"},
				Value:   "127.0.0.1:9000",
			},
			&cli.StringFlag{
				Name:    "ip",
				Aliases: []string{"i"},
				EnvVars: []string{"IP"},
				Value:   "[::]",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				EnvVars: []string{"CYBERKV_PORT"},
				Value:   5796,
			},
			&cli.IntFlag{
				Name:    "pprof",
				Aliases: []string{"pp"},
			},
		},

		Action: func(ctx *cli.Context) error {
			if ctx.Int("pprof") != 0 {
				initPprofMonitor(ctx.Int("pprof"))
			}

			endpoints := ctx.StringSlice("etcd_endpoints")
			etcdClient, err := etcdcli.New(etcdcli.Config{
				Endpoints: endpoints,
			})
			if err != nil {
				return err
			}

			minioEndpoint := ctx.String("minio_endpoint")
			minioClient, err := minio.New(minioEndpoint, &minio.Options{
				Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
				Secure: false,
			})
			c := context.Background()
			minioClient.MakeBucket(c, "cyberkv", minio.MakeBucketOptions{})
			if err != nil {
				return err
			}

			listenAddr := fmt.Sprintf("%v:%v", ctx.String("ip"), ctx.Int("port"))
			storageNode := storage.NewStorageNode(listenAddr, etcdClient, minioClient)
			storageNode.Recovery()
			storageNode.Register("storage")
			storageNode.Start()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

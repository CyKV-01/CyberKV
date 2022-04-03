package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yah01/CyberKV/common/log"
	"github.com/yah01/CyberKV/coordinator"
	etcdcli "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

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
				Name:    "ip",
				Aliases: []string{"i"},
				EnvVars: []string{"IP"},
				Value:   "[::]",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				EnvVars: []string{"CYBERKV_PORT"},
				Value:   5790,
			},
		},

		Action: func(ctx *cli.Context) error {
			endpoints := ctx.StringSlice("etcd_endpoints")
			log.Info("connecting etcd...",
				zap.Strings("endpoints", endpoints))
			etcdClient, err := etcdcli.New(etcdcli.Config{
				Endpoints: endpoints,
			})
			if err != nil {
				return err
			}
			log.Info("etcd connected")

			listenAddr := fmt.Sprintf("%v:%v", ctx.String("ip"), ctx.Int("port"))

			coord := coordinator.NewCoordinator(etcdClient, listenAddr)
			coord.Register("coordinator")
			coord.Start()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

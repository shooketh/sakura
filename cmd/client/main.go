package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/shooketh/etcd-handler"
	"github.com/shooketh/sakura"
	"github.com/shooketh/sakura/common/location"
	"github.com/shooketh/sakura/module/converter"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	endpoints     = []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"}
	timeout       = 5 * time.Second
	servicePrefix = "/sakura/service"
)

func init() {
	time.Local = location.UTC()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	handler, err := handler.New(ctx, &handler.Option{
		Config: &clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: timeout,
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := sakura.Client{
		EtcdHandler: handler,
		DialPrefix:  servicePrefix,
	}

	fmt.Println("[start] generated id")
	id, err := c.Generate(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ts, did, wid, seq := converter.Dump(id)
	fmt.Printf("id: %d\ntime: %s\ndatacenterID: %d\nworkerID: %d\nsequence: %d\n", id, ts, did, wid, seq)
	fmt.Printf("[finish] generate id\n\n")

	fmt.Println("[start] generated multiple ids")
	ids, err := c.GenerateMulti(ctx, 3)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("ids: %v\n", ids)
	fmt.Println("[finish] generate multiple ids")
}

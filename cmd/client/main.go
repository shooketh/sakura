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
	count         = 10
	endpoints     = []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"}
	timeout       = 5 * time.Second
	servicePrefix = "/sakura/service"
)

func init() {
	time.Local = location.UTC()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	handler, err := handler.New(ctx, &handler.Option{
		Config: &clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: time.Duration(timeout) * time.Second,
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("start generated id")
	start := time.Now()
	c := count

	for i := 0; i < c; i++ {
		var id int64
		var err error
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			id, err = sakura.Generate(ctx, handler, servicePrefix)
		}()
		if err != nil {
			fmt.Println(err)
			c = i
			break
		}

		ts, did, wid, seq := converter.Dump(id)
		fmt.Printf("id: %d\ntime: %s\ndatacenterID: %d\nworkerID: %d\nsequence: %d\n\n", id, ts, did, wid, seq)
	}
	fmt.Printf("time taken to generate [%d] ids is: %d ms\n", c, time.Since(start)/time.Millisecond)
	fmt.Printf("finish generate id, desired number of ids is: %d, the actual number is: %d\n", count, c)
}

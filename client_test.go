package sakura

import (
	"context"
	"testing"
	"time"

	"github.com/shooketh/etcd-handler"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	endpoints     = []string{"127.0.0.1:12379", "127.0.0.1:22379", "127.0.0.1:32379"}
	timeout       = 5 * time.Second
	servicePrefix = "/sakura/service"
)

func TestGenerate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	handler, err := handler.New(ctx, &handler.Option{
		Config: &clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: timeout,
		},
	})
	assert.Nil(t, err, "failed to create new etcd handler")

	c := Client{
		EtcdHandler: handler,
		DialPrefix:  servicePrefix,
	}
	_, err = c.Generate(ctx)
	assert.Nil(t, err, "failed to generate id")
}

func TestGenerateMulti(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	handler, err := handler.New(ctx, &handler.Option{
		Config: &clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: timeout,
		},
	})
	assert.Nil(t, err, "failed to create new etcd handler")

	c := Client{
		EtcdHandler: handler,
		DialPrefix:  servicePrefix,
	}
	_, err = c.GenerateMulti(ctx, 5)
	assert.Nil(t, err, "failed to generate multiple ids")
}

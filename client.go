package sakura

import (
	"context"
	"fmt"

	"github.com/shooketh/etcd-handler"
	"github.com/shooketh/sakura/pb"
	"google.golang.org/grpc"
)

type Client struct {
	EtcdHandler handler.Handler
	DialPrefix  string
	DialOption  []grpc.DialOption
}

func (c *Client) Generate(ctx context.Context) (int64, error) {
	conn, err := c.EtcdHandler.Dial(ctx, c.DialPrefix, c.DialOption...)
	if err != nil {
		return 0, fmt.Errorf("failed to dial etcd service: %v", err)
	}
	defer conn.Close()

	gc := pb.NewGeneratorClient(conn)

	r, err := gc.Generate(ctx, &pb.GenerateRequest{})
	if err != nil {
		return 0, fmt.Errorf("failed to generate id: %v", err)
	}

	return r.Id, nil
}

func (c *Client) GenerateMulti(ctx context.Context, n int64) ([]int64, error) {
	conn, err := c.EtcdHandler.Dial(ctx, c.DialPrefix, c.DialOption...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial etcd service: %v", err)
	}
	defer conn.Close()

	gc := pb.NewGeneratorClient(conn)

	r, err := gc.GenerateMulti(ctx, &pb.GenerateMultiRequest{Number: n})
	if err != nil {
		return nil, fmt.Errorf("failed to generate id: %v", err)
	}

	return r.Ids, nil
}

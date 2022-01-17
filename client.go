package sakura

import (
	"context"
	"fmt"

	handler "github.com/shooketh/etcd-handler"
	"github.com/shooketh/sakura/pb"
	"google.golang.org/grpc"
)

func Generate(ctx context.Context, handler handler.Handler, prefix string, opts ...grpc.DialOption) (int64, error) {
	conn, err := handler.Dial(ctx, prefix, opts...)
	if err != nil {
		return 0, fmt.Errorf("failed to dial etcd service: %v", err)
	}
	defer conn.Close()

	c := pb.NewGeneratorClient(conn)

	r, err := c.Generate(ctx, &pb.Request{})
	if err != nil {
		return 0, fmt.Errorf("failed to generate id: %v", err)
	}

	return r.Id, nil
}

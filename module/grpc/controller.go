package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/shooketh/sakura/module/generator"
	"github.com/shooketh/sakura/pb"
	"google.golang.org/grpc"
)

type Controller struct {
	IP   string
	Port string
	*generator.Generator
	net.Listener
	*grpc.Server
}

func (c *Controller) Listen(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", c.Port))
	if err != nil {
		return err
	}

	c.Listener = lis

	return nil
}

func (c *Controller) Serve(ctx context.Context) error {
	c.Server = grpc.NewServer()

	pb.RegisterGeneratorServer(c.Server, &Server{Generator: c.Generator})

	if err := c.Server.Serve(c.Listener); err != nil {
		return err
	}

	return nil
}

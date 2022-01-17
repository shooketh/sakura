package grpc

import (
	"context"

	"github.com/shooketh/sakura/common/log"
	"github.com/shooketh/sakura/module/generator"
	"github.com/shooketh/sakura/pb"
)

type Server struct {
	pb.UnimplementedGeneratorServer
	*generator.Generator
}

// Generate implements sakura.GeneratorServer
func (s *Server) Generate(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	id := s.Generator.Generate()

	log.Logger.Debug().Msgf("generated ID: %d", id)

	return &pb.Reply{Id: id}, nil
}

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
func (s *Server) Generate(ctx context.Context, in *pb.GenerateRequest) (*pb.GenerateReply, error) {
	id := s.Generator.Generate()

	log.Logger.Debug().Msgf("generated id: %d", id)

	return &pb.GenerateReply{Id: id}, nil
}

func (s *Server) GenerateMulti(ctx context.Context, in *pb.GenerateMultiRequest) (*pb.GenerateMultiReply, error) {
	var ids []int64
	for i := 0; i < int(in.Number); i++ {
		ids = append(ids, s.Generator.Generate())
	}

	log.Logger.Debug().Msgf("generated ids: %v", ids)

	return &pb.GenerateMultiReply{Ids: ids}, nil
}

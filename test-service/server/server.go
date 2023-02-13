package server

import (
	"context"
	pb "test-service/gen/go/proto"
)

type Service struct {
	*pb.UnimplementedGreeterServiceServer
}

func (s *Service) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{Message: "Hello " + req.Name}, nil
}

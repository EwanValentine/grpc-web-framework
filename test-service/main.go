package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "test-service/gen/go/proto"
	"test-service/server"
)

func main() {
	svc := &server.Service{}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	pb.RegisterGreeterServiceServer(s, svc)

	if err := s.Serve(lis); err != nil {
		log.Panic(err)
	}
}

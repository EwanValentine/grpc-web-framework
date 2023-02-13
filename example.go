package main

import (
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"grpc-gateway/server"
	"log"
	"net/http"
	pb "test-service/gen/go/proto"
)

func main() {
	lis, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}

	client := pb.NewGreeterServiceClient(lis)
	s := server.New()

	s.RegisterEndpoint(http.MethodPost, "/", server.NewHandler[*pb.GreetRequest, *pb.GreetResponse](client.Greet, func(i []byte) (*pb.GreetRequest, error) {
		r := new(pb.GreetRequest)
		if err := protojson.Unmarshal(i, r); err != nil {
			return r, err
		}
		return r, nil
	}))

	log.Fatal(http.ListenAndServe(":8081", s))
}

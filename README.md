# GRPC Web Proxy Framework

*Status: very much a work in progress*

## Description

A Go framework for building a GRPC web proxy. It's a bit like Envoy, but hopefully, much simpler.

This differs from grpc-web, because it's intended to proxy multiple services into one web service.

## Motivation

- GRPC-web is great if you want to generate a webserver for a GRPC service. However, it's not that useful if you want to proxy several services.
- Envoy is a bit of a pain to set up and configure. So I wanted something simpler, that could be configured via Go code.
- I couldn't find any other projects that did what I wanted, so, decided to try to make it myself.

## Example

```golang
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
```

As you can see, we use generics to configure the type of the request and response, and takes the client method and a callback. 

The callback is a 'Marshaller', which marshalls the request body into the request type.

## TODO
- [ ] Add support for streaming
- [ ] Add support for mapping headers into context metadata
- [ ] Add support for converting get request params into a GRPC type

package main

import (
	"context"
	pb "github.com/koala/example/grpc_example/hello"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Reply: "你好 " + in.Name,
	}, nil
}

func main() {

}

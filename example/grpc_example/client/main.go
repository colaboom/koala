package main

import (
	pb "github.com/koala/example/grpc_example/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const (
	address     = "localhost:12345"
	defaultName = "colaboom"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect, err :%v", err)
	}
	defer conn.Close()
	c := pb.NewHelloServiceClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	for {
		resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatal("could not greet :%v", err)
		}
		log.Printf("Greeting: %s", resp.Reply)
		time.Sleep(time.Millisecond * 10)
	}
}

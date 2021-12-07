package main

import (
	"context"
	"fmt"
	"github.com/koala/example/grpc_example/hello"
	"github.com/koala/logs"
	"github.com/koala/rpc"
	"google.golang.org/grpc"
	"time"
)

type HelloClient struct {
	serviceName string
}

func NewHelloClient(serviceName string) *HelloClient {
	return &HelloClient{
		serviceName: serviceName,
	}
}

func (h *HelloClient) SayHelloV1(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(context.Background(), "did not connect, err :%v", err)
		return nil, err
	}
	defer conn.Close()

	c := hello.NewHelloServiceClient(conn)
	resp, err := c.SayHello(ctx, in, opts...)
	if err != nil {
		logs.Error(ctx, "could not greet :%v", err)
		return nil, err
	}
	logs.Info(ctx, "Greeting: %s", resp.Reply)

	return resp, err
}

func (h *HelloClient) SayHelloV2(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error) {
	middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
	mkResp, err := middlewareFunc(ctx, in)
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*hello.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *hello.HelloResponse")
	}

	return resp, err
}

func mwClientSayHello(ctx context.Context, request interface{}) (response interface{}, err error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect, err :%v", err)
		return nil, err
	}
	defer conn.Close()

	req := request.(*hello.HelloRequest)
	c := hello.NewHelloServiceClient(conn)
	return c.SayHello(ctx, req)
}

func myClientExample() {
	client := NewHelloClient("hello")
	ctx := context.Background()
	in := &hello.HelloRequest{Name: "client v2"}
	resp, err := client.SayHelloV2(ctx, in)
	if err != nil {
		logs.Error(ctx, "could not greet :%v", err)
		return
	}

	logs.Info(ctx, "Greeting: %s", resp.Reply)
	fmt.Printf("Greeting: %s", resp.Reply)
	return
}

func main() {
	myClientExample()
	time.Sleep(time.Millisecond * 100)
}

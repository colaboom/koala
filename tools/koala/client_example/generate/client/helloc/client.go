package helloc

import (
	"context"
	"fmt"
	"github.com/koala/tools/koala/client_example/generate/hello"
	"github.com/koala/logs"
	"github.com/koala/rpc"
	"google.golang.org/grpc"
)

type HelloClient struct {
	serviceName string
	client *rpc.KoalaClient
}

func NewHelloClient(serviceName string, opts...rpc.RpcOptionFunc) *HelloClient {
	c := &HelloClient{
		serviceName: serviceName,
	}
	c.client = rpc.NewKoalaClient(serviceName, opts...)
	return c
}


func (h *HelloClient) SayHello(ctx context.Context, in *hello.HelloRequest, opts ...grpc.CallOption) (*hello.HelloResponse, error) {
	/*middlewareFunc := rpc.BuildClientMiddleware(mwClientSayHello)
	mkResp, err := middlewareFunc(ctx, in)
	if err != nil {
		return nil, err
	}*/

	mkResp, err := h.client.Call(ctx, "SayHello", in, mwClientSayHello)
	resp, ok := mkResp.(*hello.HelloResponse)
	if !ok {
		err = fmt.Errorf("invalid resp, not *hello.HelloResponse")
	}

	return resp, err
}

func mwClientSayHello(ctx context.Context, request interface{}) (response interface{}, err error) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect, err :%v", err)
		return nil, err
	}
	defer conn.Close()

	req := request.(*hello.HelloRequest)
	c := hello.NewHelloServiceClient(conn)
	return c.SayHello(ctx, req)
}


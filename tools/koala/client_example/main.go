package main

import (
	"context"
	"fmt"
	"github.com/koala/logs"
	"github.com/koala/tools/koala/client_example/generate/client/helloc"
	"github.com/koala/tools/koala/client_example/generate/hello"
	"time"
)

func myClientExample() {
	client := helloc.NewHelloClient("hello")
	ctx := context.Background()
	in := &hello.HelloRequest{Name: "client"}
	resp, err := client.SayHello(ctx, in)
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

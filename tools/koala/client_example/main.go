package main

import (
	"context"
	"github.com/koala/logs"
	"github.com/koala/rpc"
	"github.com/koala/tools/koala/client_example/generate/client/helloc"
	"github.com/koala/tools/koala/client_example/generate/hello"
	"time"
)

func myClientExample() {
	client := helloc.NewHelloClient("hello", rpc.WithLimitQPS(5))
	var count int
	for {
		count++
		ctx := context.Background()
		resp, err := client.SayHello(ctx, &hello.HelloRequest{Name: "test my client"})
		if err != nil {
			if count%100 == 0 {
				logs.Error(ctx, "could not greet: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
			continue
		}

		logs.Info(ctx, "Greeting: %s", resp.Reply)

		time.Sleep(100 * time.Millisecond)
	}
	/*ctx := context.Background()
	in := &hello.HelloRequest{Name: "client"}
	resp, err := client.SayHello(ctx, in)
	if err != nil {
		logs.Error(ctx, "could not greet :%v", err)
		//logs.Stop()
		return
	}

	logs.Info(ctx, "Greeting: %s", resp.Reply)
	//logs.Stop()
	fmt.Printf("Greeting: %s\n", resp.Reply)
	return*/
}

func main() {
	myClientExample()
	time.Sleep(time.Millisecond * 100)
}

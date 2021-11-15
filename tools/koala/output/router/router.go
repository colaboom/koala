
package router

import (
	"golang.org/x/net/context"
	
	"github.com/koala/tools/koala/output/generate/hello"
	

	
	"github.com/koala/tools/koala/output/controller"
	
)

type RouterServer struct{}


func(s *RouterServer) SayHello(ctx context.Context, r *hello.HelloRequest)(resp *hello.HelloResponse, err error){
	ctrl :=&controller.SayHelloController{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}


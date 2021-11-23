
package router

import (
	"golang.org/x/net/context"
	//"github.com/koala/middleware"
	"github.com/koala/meta"
	
	"github.com/koala/tools/koala/output/generate/hello"
	

	
	"github.com/koala/tools/koala/output/controller"
	
)

type RouterServer struct{}


func(s *RouterServer) SayHello(ctx context.Context, r *hello.HelloRequest)(resp *hello.HelloResponse, err error){
	ctx = meta.InitServerMeta(ctx, "hello", "SayHello")
	/*mwFunc := middleware.BuildServerMiddleware(mwSayHello)
	mwResp, err := mwFunc(ctx, r)

	resp = mwResp.(*hello.HelloResponse)*/
	respTmp, err := mwSayHello(ctx, r)
	resp = respTmp.(*hello.HelloResponse)
	return
}

func mwSayHello(ctx context.Context, request interface{}) (resp interface{}, err error) {
	r := request.(*hello.HelloRequest)
	ctrl := &controller.SayHelloController{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}



package controller

import(
	"golang.org/x/net/context"
	
	"github.com/koala/tools/koala/output/generate/hello"
	
)

type SayHelloController struct{}
//检查参数
func (s *SayHelloController) CheckParams(ctx context.Context, r *hello.HelloRequest) (err error) {
	return
}

//SayHello函数的实现
func (s *SayHelloController) Run(ctx context.Context, r *hello.HelloRequest) (resp *hello.HelloResponse, err error){
	return
}

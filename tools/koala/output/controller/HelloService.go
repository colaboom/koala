package controller

import(
"golang.org/x/net/context"
"github.com/koala/tools/koala/output/generate"
)
type Server struct{}

func (s *Server) SayHello(ctx context.Context, r *hello.HelloRequest)(resp *hello.HelloResponse, err error){
return
}

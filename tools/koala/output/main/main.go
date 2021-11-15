
package main

import(
	"net"
	"log"
	"google.golang.org/grpc"
	
	"github.com/koala/tools/koala/output/router"
	
	
	"github.com/koala/tools/koala/output/generate/hello"
	
)
var server = &router.RouterServer{}
var port = ":12345"

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen :%v", err)
	}
	
	s := grpc.NewServer()
	hello.RegisterHelloServiceServer(s, server)
	s.Serve(listen)
}

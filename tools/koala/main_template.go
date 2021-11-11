package main

var main_template = `
package main

import(
	"net"
	"log"
	"google.golang.org/grpc"
	"github.com/koala/tools/koala/output/controller"
	hello "github.com/koala/tools/koala/output/generate"
)
var server = &controller.Server{}
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
`

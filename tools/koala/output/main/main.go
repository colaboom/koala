package main

import(
"net"
"log"
"google.golang.org/grpc"
"github.com/koala/tools/koala/output/controller"
pb "github.com/koala/tools/koala/output/generate"
)
var server = &controller.Server{}

var port = ":12345"


func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen :%!v(MISSING)", err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, server)
	s.Serve(listen)
}


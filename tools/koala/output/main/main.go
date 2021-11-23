
package main

import(
	"net"
	"log"
	"google.golang.org/grpc"
    "github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	
	"github.com/koala/tools/koala/output/router"
	
	
	"github.com/koala/tools/koala/output/generate/hello"
	
)
var server = &router.RouterServer{}
var port = ":12345"

func main() {

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe("0.0.0.0:9091", nil))
	}()

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen :%v", err)
	}
	
	s := grpc.NewServer()
	hello.RegisterHelloServiceServer(s, server)
	s.Serve(listen)
}

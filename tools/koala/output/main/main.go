
package main

import(
	"net"
	"log"
	"fmt"
	"google.golang.org/grpc"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/koala/server"
	"net/http"
	
	"github.com/koala/tools/koala/output/router"
	
	
	"github.com/koala/tools/koala/output/generate/hello"
	
)
var routerServer = &router.RouterServer{}

func main() {

	err := server.Init("hello")
	if err != nil {
		log.Fatal("init server failed, err :%v\n", err)
	}

	if server.GetConf().Prometheus.SwitchOn {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			addr := fmt.Sprintf("0.0.0.0:%d", server.GetConf().Prometheus.Port)
			log.Fatal(http.ListenAndServe(addr, nil))
		}()
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", server.GetConf().Port))
	if err != nil {
		log.Fatal("failed to listen :%v", err)
	}
	
	s := grpc.NewServer()
	hello.RegisterHelloServiceServer(s, routerServer)
	s.Serve(listen)
}

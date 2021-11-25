
package main

import(
	"log"
    "github.com/koala/server"
	
	"github.com/koala/tools/koala/output/router"
	
	
	"github.com/koala/tools/koala/output/generate/hello"
	
)
var routerServer = &router.RouterServer{}

func main() {

	err := server.Init("hello")
	if err != nil {
		log.Fatal("init server failed, err :%v\n", err)
	}

	hello.RegisterHelloServiceServer(server.GRPCServer(), routerServer)
	server.Run()
}

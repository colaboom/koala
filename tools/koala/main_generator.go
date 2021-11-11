package main

import (
	"fmt"
	"os"
	"path"
)

type MainGenerator struct {
}

func (d *MainGenerator) Run(opt *Option) (err error) {
	filename := path.Join(opt.Output, "main/main.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file %s failed, err :%v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "package main\n")
	fmt.Fprintln(file)
	fmt.Fprintf(file, "import(\n")
	fmt.Fprintf(file, `"net"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `"log"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `"google.golang.org/grpc"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `"github.com/koala/tools/koala/output/controller"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `pb "github.com/koala/tools/koala/output/generate"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, ")\n")
	fmt.Fprintf(file, "var server = &controller.Server{}\n")
	fmt.Fprintf(file, "\n")
	fmt.Fprintf(file, "var port = \":12345\"\n")
	fmt.Fprintf(file, "\n")
	fmt.Fprintf(file,`
func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen :%v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, server)
	s.Serve(listen)
}
`)
	fmt.Fprintf(file, "\n")

	return
}

func init() {
	dir := &MainGenerator{}

	Register("main generator", dir)
}

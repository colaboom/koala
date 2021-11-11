package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"path"
)

type CtrlGenerator struct {
	service *proto.Service
	message []*proto.Message
	rpc     []*proto.RPC
}

func (d *CtrlGenerator) Run(opt *Option) (err error) {
	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Printf("open file %s failed,err :%v\n", opt.Proto3Filename, err)
		return
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		fmt.Printf("parse file %s failed,err:%v\n", opt.Proto3Filename, err)
	}

	proto.Walk(definition,
		proto.WithService(d.handleService),
		proto.WithMessage(d.handleMessage),
		proto.WithRPC(d.handleRPC))

	//fmt.Printf("parse protoc succ, rpc:%#v\n", d.rpc)
	return d.generateRpc(opt)
}

func (d *CtrlGenerator) generateRpc(opt *Option) (err error) {
	filename := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", d.service.Name))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file %s failed, err : %v\n", filename, err)
	}
	defer file.Close()

	fmt.Fprintf(file, "package controller\n")
	fmt.Fprintln(file)
	fmt.Fprintf(file, "import(\n")
	fmt.Fprintf(file, `"golang.org/x/net/context"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, `"github.com/koala/tools/koala/output/generate"`)
	fmt.Fprintln(file)
	fmt.Fprintf(file, ")\n")
	fmt.Fprintf(file, "type Server struct{}\n")
	fmt.Fprintf(file, "\n")
	for _, rpc := range d.rpc {
		fmt.Fprintf(file, "func (s *Server) %s(ctx context.Context, r *hello.%s)(resp *hello.%s, err error){\nreturn\n}\n", rpc.Name, rpc.RequestType, rpc.ReturnsType)
	}
	return
}

func (d *CtrlGenerator) handleService(s *proto.Service) {
	//fmt.Println(s.Name)
	d.service = s
}

func (d *CtrlGenerator) handleMessage(m *proto.Message) {
	//fmt.Println(m.Name)
	d.message = append(d.message, m)
}

func (d *CtrlGenerator) handleRPC(r *proto.RPC) {
	/*fmt.Println(r.Name)
	fmt.Println(r.RequestType)
	fmt.Println(r.ReturnsType)
	fmt.Printf("rpc:%#v, comment:%v\n", r, r.Comment)*/
	d.rpc = append(d.rpc, r)
}

func init() {
	dir := &CtrlGenerator{}

	Register("ctrl generator", dir)
}

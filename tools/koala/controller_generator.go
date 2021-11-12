package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"path"
	"text/template"
)

type CtrlGenerator struct {
	/*service *proto.Service
	message []*proto.Message
	rpc     []*proto.RPC*/
}

type RpcMeta struct {
	Rpc     *proto.RPC
	Package *proto.Package
	Prefix  string
}

func (d *CtrlGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	reader, err := os.Open(opt.Proto3Filename)
	if err != nil {
		fmt.Printf("open file %s failed,err :%v\n", opt.Proto3Filename, err)
		return
	}
	defer reader.Close()

	return d.generateRpc(opt, metaData)
}

func (d *CtrlGenerator) generateRpc(opt *Option, metaData *ServiceMetaData) (err error) {
	for _, rpc := range metaData.Rpc {
		var file *os.File
		filename := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", rpc.Name))
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Printf("open file %s failed, err : %v\n", filename, err)
		}
		rpcMeta := &RpcMeta{}
		rpcMeta.Rpc = rpc
		rpcMeta.Package = metaData.Package
		rpcMeta.Prefix = metaData.Prefix
		err = d.render(file, controller_template, rpcMeta)
		if err != nil {
			fmt.Printf("render controller failed,err :%v\n", err)
			return
		}
		defer file.Close()
	}
	return
}

func (d *CtrlGenerator) render(file *os.File, data string, rpcMeta *RpcMeta) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		return
	}
	err = t.Execute(file, rpcMeta)
	if err != nil {
		return
	}
	return
}

func init() {
	dir := &CtrlGenerator{}

	Register("ctrl generator", dir)
}

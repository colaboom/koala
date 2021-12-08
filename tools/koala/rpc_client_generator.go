package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

type RpcClientGenerator struct {

}

func (d *RpcClientGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	//所有生成的rpc client的包名都以c后缀结尾，最后一级目录加上c的后缀
	packageName := fmt.Sprintf("%sc", metaData.Package.Name)
	dir := path.Join(opt.Output, "generate", "client", packageName)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("mkdir dir %s failed, err : %v\n", dir, err)
		return
	}
	filename := path.Join(dir, "client.go")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("open file %s failed, err :%v\n", filename, err)
		return
	}
	defer file.Close()
	err = d.render(file, rpcClientTemplate, metaData)
	if err != nil {
		fmt.Printf("render failed, err :%v\n", err)
		return
	}

	return
}

func (d *RpcClientGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error){
	t := template.New("main").Funcs(templateFuncMap)
	t, err = t.Parse(data)
	if err != nil {
		return
	}
	err = t.Execute(file, metaData)
	if err != nil {
		return
	}
	return
}

func init() {
	dir := &RpcClientGenerator{}

	RegisterClientGenerator("rpc client generator", dir)
}
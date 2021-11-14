package main

var controller_template = `
package controller

import(
	"golang.org/x/net/context"
	{{if not .Prefix}}
	"generate/ {{.Package.Name}}"
	{{else}}
	"{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}
)

type {{.Rpc.Name}}Controller struct{}
//检查参数
func (s *{{.Rpc.Name}}Controller) CheckParams(ctx context.Context, r *{{$.Package.Name}}.{{.Rpc.RequestType}}) (err error) {
	return
}

//SayHello函数的实现
func (s *{{.Rpc.Name}}Controller) Run(ctx context.Context, r *{{$.Package.Name}}.{{.Rpc.RequestType}}) (resp *{{$.Package.Name}}.{{.Rpc.ReturnsType}}, err error){
	return
}
`

package main

var router_template = `
package router

import (
	"golang.org/x/net/context"
	{{if not .Prefix}}
	"{{.Package.Name}}"
	{{else}}
	{{.Package.Name}} "{{.Prefix}}/generate"
	{{end}}

	{{if not .Prefix}}
	"controller"
	{{else}}
	"{{.Prefix}}/controller"
	{{end}}
)

type RouterServer struct{}

{{range .Rpc}}
func(s *RouterServer) {{.Name}}(ctx context.Context, r *{{$.Package.Name}}.{{.RequestType}})(resp *{{$.Package.Name}}.{{.ReturnsType}})
	ctrl :=&controller.{{.Name}}Controller{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}
{{end}}
`

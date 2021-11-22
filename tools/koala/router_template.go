package main

var router_template = `
package router

import (
	"golang.org/x/net/context"
	{{if not .Prefix}}
	"generate/{{.Package.Name}}"
	{{else}}
	"{{.Prefix}}/generate/{{.Package.Name}}"
	{{end}}

	{{if not .Prefix}}
	"controller"
	{{else}}
	"{{.Prefix}}/controller"
	{{end}}
)

type RouterServer struct{}

{{range .Rpc}}
func(s *RouterServer) {{.Name}}(ctx context.Context, r *{{$.Package.Name}}.{{.RequestType}})(resp *{{$.Package.Name}}.{{.ReturnsType}}, err error){
	ctx = meta.InitServerMeta(ctx, "{{$.Package.Name}}", "{{.Name}}")
	mwFunc := middleware.BuildServerMiddleware(mw{{.Name}})
	mwResp, err := mwFunc(ctx, r)

	resp = mwResp.(*{{$.Package.Name}}.{{.ReturnsType}})
	return
}

func mw{{.Name}}(ctx context.Context, request interface{}) (resp interface{}, err error) {
	r := request.(*{{$.Package.Name}}.{{.RequestType}})
	ctrl := &controller.{{.Name}}Controller{}
	err = ctrl.CheckParams(ctx, r)
	if err != nil {
		return
	}

	resp, err = ctrl.Run(ctx, r)
	return
}
{{end}}
`

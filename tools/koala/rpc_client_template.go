package main

var rpcClientTemplate = `package {{.Package.Name}}c

import (
	"context"
	"fmt"
	"{{.Prefix}}/generate/{{.Package.Name}}"
	"github.com/koala/logs"
	"github.com/koala/rpc"
	"google.golang.org/grpc"
)

type {{Capitalize .Package.Name}}Client struct {
	serviceName string
}

func New{{Capitalize .Package.Name}}Client(serviceName string) *{{Capitalize .Package.Name}}Client {
	return &{{Capitalize .Package.Name}}Client{
		serviceName: serviceName,
	}
}

{{range .Rpc}}
func (h *{{Capitalize $.Package.Name}}Client) {{.Name}}(ctx context.Context, in *{{$.Package.Name}}.{{.RequestType}}, opts ...grpc.CallOption) (*{{$.Package.Name}}.{{.ReturnsType}}, error) {
	middlewareFunc := rpc.BuildClientMiddleware(mwClient{{.Name}})
	mkResp, err := middlewareFunc(ctx, in)
	if err != nil {
		return nil, err
	}

	resp, ok := mkResp.(*{{$.Package.Name}}.{{.ReturnsType}})
	if !ok {
		err = fmt.Errorf("invalid resp, not *{{$.Package.Name}}.{{.ReturnsType}}")
	}

	return resp, err
}

func mwClient{{.Name}}(ctx context.Context, request interface{}) (response interface{}, err error) {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		logs.Error(ctx, "did not connect, err :%v", err)
		return nil, err
	}
	defer conn.Close()

	req := request.(*{{$.Package.Name}}.{{.RequestType}})
	c := {{$.Package.Name}}.New{{Capitalize $.Package.Name}}ServiceClient(conn)
	return c.{{.Name}}(ctx, req)
}
{{end}}
`

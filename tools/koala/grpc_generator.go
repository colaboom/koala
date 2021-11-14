package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type GrpcGenerator struct {
	dirList []string
}

func (d *GrpcGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	// protoc --go_out=plugins=grpc:. hello.proto
	dir := path.Join(opt.Output, "generate", metaData.Package.Name)
	os.MkdirAll(dir, 0755)
	outputParams := fmt.Sprintf("plugins=grpc:%s/generate/%s", opt.Output, metaData.Package.Name)
	cmd := exec.Command("protoc", "--go_out", outputParams, opt.Proto3Filename)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Printf("grpc generate failed, err :%v\n", err)
		return
	}
	return
}

func init() {
	dir := &GrpcGenerator{}

	Register("grpc generator", dir)
}

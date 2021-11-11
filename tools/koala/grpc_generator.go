package main

import (
	"fmt"
	"os"
	"os/exec"
)

type GrpcGenerator struct {
	dirList []string
}

func (d *GrpcGenerator) Run(opt *Option) (err error) {
	// protoc --go_out=plugins=grpc:. hello.proto
	outputParams := fmt.Sprintf("plugins=grpc:%s/generate/", opt.Output)
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

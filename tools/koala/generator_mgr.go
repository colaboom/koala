package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"path"
)

var genMgr = &GeneratorMgr{
	genMap:   make(map[string]Generator),
	metaData: &ServiceMetaData{},
}

type GeneratorMgr struct {
	genMap   map[string]Generator
	metaData *ServiceMetaData
}

var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"scripts",
	"conf",
	"app/router",
	"app/config",
	"model",
	"generate",
}

func (g *GeneratorMgr) parseService(opt *Option) (err error) {
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
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRPC))

	return
}

func (g *GeneratorMgr) handleService(s *proto.Service) {
	//fmt.Println(s.Name)
	g.metaData.Service = s
}

func (g *GeneratorMgr) handleMessage(m *proto.Message) {
	//fmt.Println(m.Name)
	g.metaData.Message = append(g.metaData.Message, m)
}

func (g *GeneratorMgr) handleRPC(r *proto.RPC) {
	/*fmt.Println(r.Name)
	fmt.Println(r.RequestType)
	fmt.Println(r.ReturnsType)
	fmt.Printf("rpc:%#v, comment:%v\n", r, r.Comment)*/
	g.metaData.Rpc = append(g.metaData.Rpc, r)
}

func (g *GeneratorMgr) createAllDir(opt *Option) (err error) {
	for _, dir := range AllDirList {
		fullDir := path.Join(opt.Output, dir)
		err = os.MkdirAll(fullDir, 0755)
		if err != nil {
			fmt.Printf("mkdir dir %s failed, err : %v\n", dir, err)
			return
		}
	}
	return
}

func (g *GeneratorMgr) Run(opt *Option) (err error) {
	err = g.parseService(opt)
	if err != nil {
		return
	}
	err = g.createAllDir(opt)
	if err != nil {
		return
	}
	for _, gen := range g.genMap {
		err = gen.Run(opt, g.metaData)
		if err != nil {
			return
		}
	}
	return
}

func Register(name string, gen Generator) (err error) {
	_, ok := genMgr.genMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists", name)
	}

	genMgr.genMap[name] = gen
	return
}

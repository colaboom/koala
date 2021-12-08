package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var genMgr = &GeneratorMgr{
	genServerMap: make(map[string]Generator),
	genClientMap: make(map[string]Generator),
	metaData:     &ServiceMetaData{},
}

type GeneratorMgr struct {
	genServerMap map[string]Generator
	genClientMap map[string]Generator
	metaData     *ServiceMetaData
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
	"router",
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
		proto.WithRPC(g.handleRPC),
		proto.WithPackage(g.handlePackage))

	return
}

func (g *GeneratorMgr) handleService(s *proto.Service) {
	//fmt.Println("serviceName:", s.Name)
	g.metaData.Service = s
}

func (g *GeneratorMgr) handleMessage(m *proto.Message) {
	//fmt.Println("messageName:", m.Name)
	g.metaData.Message = append(g.metaData.Message, m)
}

func (g *GeneratorMgr) handleRPC(r *proto.RPC) {
	//fmt.Println("rpcName", r.Name)
	//fmt.Println("rpcRequestType", r.RequestType)
	//fmt.Println("rpcReturnType", r.ReturnsType)
	//fmt.Printf("rpc:%#v, comment:%v\n", r, r.Comment)
	g.metaData.Rpc = append(g.metaData.Rpc, r)
}

func (g *GeneratorMgr) handlePackage(p *proto.Package) {
	//fmt.Println("packageName", p.Name)
	g.metaData.Package = p
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

func (g *GeneratorMgr) initOutputDir(opt *Option) /* (err error)*/ {
	goPath := os.Getenv("GOPATH")
	// 指定路径
	if len(opt.Prefix) > 0 {
		opt.Output = path.Join(goPath, "src", opt.Prefix)
		return
	}

	// 没有指定路径，就用当前路径
	exeFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}

	if runtime.GOOS == "windows" {
		exeFilePath = strings.Replace(exeFilePath, "\\", "/", -1)
	}

	lastIdx := strings.LastIndex(exeFilePath, "/")
	if lastIdx < 0 {
		err = fmt.Errorf("invalid exe exeFilePath:%v", exeFilePath)
		return
	}

	tmpGoPath := strings.ToLower(goPath)
	tmpGoPath = strings.Replace(tmpGoPath, "\\", "/", -1)
	opt.Output = strings.ToLower(exeFilePath[0:lastIdx])
	srcPath := path.Join(tmpGoPath, "src/")
	if srcPath[len(srcPath)-1] != '/' {
		srcPath = fmt.Sprintf("%s/", srcPath)
	}
	opt.Prefix = strings.Replace(opt.Output, srcPath, "", -1)
	return
}

func (g *GeneratorMgr) Run(opt *Option) (err error) {
	err = g.parseService(opt)
	if err != nil {
		return
	}
	g.initOutputDir(opt)
	g.metaData.Prefix = opt.Prefix

	if opt.GenServerCode {
		err = g.createAllDir(opt)
		if err != nil {
			return
		}
		for _, gen := range g.genServerMap {
			err = gen.Run(opt, g.metaData)
			if err != nil {
				return
			}
		}
	}

	if opt.GenClientCode {
		for _, gen := range g.genClientMap {
			err = gen.Run(opt, g.metaData)
			if err != nil {
				return
			}
		}
	}
	return
}

func RegisterServerGenerator(name string, gen Generator) (err error) {
	_, ok := genMgr.genServerMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists", name)
	}

	genMgr.genServerMap[name] = gen
	return
}

func RegisterClientGenerator(name string, gen Generator) (err error) {
	_, ok := genMgr.genClientMap[name]
	if ok {
		err = fmt.Errorf("generator %s is exists", name)
	}

	genMgr.genClientMap[name] = gen
	return
}

package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	var opt Option

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Value:       "./test.proto",
			Usage:       "idl filename",
			Destination: &opt.Proto3Filename,
		},
		cli.BoolFlag{
			Name:        "c",
			Usage:       "generate grpc client code",
			Destination: &opt.GenClientCode,
		},
		cli.BoolFlag{
			Name:        "s",
			Usage:       "generate grpc server code",
			Destination: &opt.GenServerCode,
		},
		cli.StringFlag{
			Name:        "p",
			Value:       "",
			Usage:       "prefix of package",
			Destination: &opt.Prefix,
		},
	}

	app.Action = func(c *cli.Context) error {
		// 命令行程序代码的入口
		err := genMgr.Run(&opt)
		if err != nil {
			fmt.Printf("code generate failed, err :%v", err)
			return err
		}
		fmt.Printf("code generate succ")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

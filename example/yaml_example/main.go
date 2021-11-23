package main

import (
	"fmt"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Site  SiteConfig  `yaml:"site"`
	Nginx NginxConfig `yaml:"nginx"`
}

type SiteConfig struct {
	Port      int    `yaml:"port"`
	HttpsOn   bool   `yaml:"https_on"`
	Domain    string `yaml:"domain"`
	HttpsPort int    `yaml:"https_port"`
}

type NginxConfig struct {
	Port     int      `yaml:"port"`
	LogPath  string   `yaml:"log_path"`
	Path     string   `yaml:"path"`
	SiteName string   `yaml:"site_name"`
	SiteAddr string   `yaml:"site_addr"`
	Upstream []string `yaml:"upstream"`
}

func main() {
	file, err := ioutil.ReadFile("./test.yaml")
	if err != nil {
		fmt.Printf("init yaml file failed, err :%v\n", err)
		return
	}

	var conf Config
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		fmt.Printf("unmarshal file failed, err :%v\n", err)
		return
	}

	fmt.Printf("site port:%d\n", conf.Site.Port)
	fmt.Printf("conf:%#v\n", conf)
}

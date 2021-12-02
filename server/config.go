package server

import (
	"fmt"
	"github.com/koala/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	// 默认配置
	koalaConf = &KoalaConf{
		Port: 8080,
		Prometheus: PrometheusConf{
			SwitchOn: true,
			Port:     8081,
		},
		ServiceName: "koala_server",
		Register: RegisterConf{
			SwitchOn: false,
		},
		Log: LogConf{
			Level: "debug",
			Dir:   "./logs/",
		},
		Limit: LimitConf{
			SwitchOn: true,
			QPSLimit: 50000,
		},
	}
)

type KoalaConf struct {
	Port        int            `yaml:"port"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	ServiceName string         `yaml:"service_name"`
	Register     RegisterConf   `yaml:"register"`
	Log         LogConf        `yaml:"log"`
	Limit       LimitConf      `yaml:"limit"`

	//内部的配置项
	ConfigDir  string `yaml:"-"`
	RootDir    string `yaml:"-"`
	ConfigFile string `yaml:"-"`
}

type LimitConf struct {
	QPSLimit int  `yaml:"qps"`
	SwitchOn bool `yaml:"switch_on"`
}

type PrometheusConf struct {
	SwitchOn bool `yaml:"switch_on"`
	Port     int  `yaml:"port"`
}

type LogConf struct {
	Level      string `yaml:"level"`
	Dir        string `yaml:"path"`
	ChanSize   int    `yaml:"chan_size"`
	ConsoleLog bool   `yaml:"console_log"`
}

type RegisterConf struct {
	SwitchOn     bool          `yaml:"switch_on"`
	RegisterPath string        `yaml:"register_path"`
	Timeout      time.Duration `yaml:"timeout"`
	HeartBeat    int64         `yaml:"heart_beat"`
	RegisterName string        `yaml:"register_name"`
	RegisterAddr string        `yaml:"register_addr"`
}

func initDir(serverName string) (err error) {
	exeFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}

	if runtime.GOOS == "windows" {
		exeFilePath = strings.Replace(exeFilePath, "\\", "/", -1)
	}

	lastIndex := strings.LastIndex(exeFilePath, "/")

	if lastIndex < 0 {
		err = fmt.Errorf("invalid exe path:%v\n", exeFilePath)
		return
	}

	koalaConf.RootDir = path.Join(strings.ToLower(exeFilePath[0:lastIndex]), "..")
	koalaConf.ConfigDir = path.Join(koalaConf.RootDir, "conf", util.GetEnv())
	koalaConf.ConfigFile = path.Join(koalaConf.ConfigDir, fmt.Sprintf("%s.yaml", serverName))

	return
}

func InitConfig(serverName string) (err error) {
	err = initDir(serverName)
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(koalaConf.ConfigFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &koalaConf)
	if err != nil {
		return
	}

	fmt.Printf("init koala config succ, conf:%#v\n", koalaConf)
	return
}

func GetConfigDir() string {
	return koalaConf.ConfigDir
}

func GetRootDir() string {
	return koalaConf.RootDir
}

func GetServerPort() int {
	return koalaConf.Port
}

func GetConf() *KoalaConf {
	return koalaConf
}

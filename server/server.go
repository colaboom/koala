package server

import (
	"context"
	"fmt"
	"github.com/koala/logs"
	"github.com/koala/middleware"
	"github.com/koala/registry"
	_ "github.com/koala/registry/etcd"
	"github.com/koala/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type KoalaServer struct {
	*grpc.Server
	limiter  *rate.Limiter
	register registry.Registry

	userMiddleware []middleware.Middleware
}

var koalaServer = &KoalaServer{
	Server: grpc.NewServer(),
}

func Use(m ...middleware.Middleware) {
	koalaServer.userMiddleware = append(koalaServer.userMiddleware, m...)
}

func Init(serviceName string) (err error) {
	err = InitConfig(serviceName)
	if err != nil {
		return
	}

	// 初始化限流器
	if koalaConf.Limit.SwitchOn {
		koalaServer.limiter = rate.NewLimiter(rate.Limit(koalaConf.Limit.QPSLimit), koalaConf.Limit.QPSLimit)
	}

	// 初始化日志
	initLogger()

	// 初始化注册中心
	err = initRegister(serviceName)
	if err != nil {
		logs.Error(context.TODO(), "init register failed , err :%v", err)
		return
	}

	return
}

func initLogger() (err error) {
	filename := fmt.Sprintf("%s/%s.log", koalaConf.Log.Dir, koalaConf.ServiceName)
	outputer, err := logs.NewFileOutputer(filename)
	if err != nil {
		return
	}

	level := logs.GetLogLevel(koalaConf.Log.Level)
	logs.InitLogger(level, koalaConf.Log.ChanSize, koalaConf.ServiceName)
	logs.AddOutputer(outputer)

	if koalaConf.Log.ConsoleLog {
		logs.AddOutputer(logs.NewConsoleOutputer())
	}
	return
}

func initRegister(serviceName string) (err error) {
	if !koalaConf.Register.SwitchOn {
		return
	}

	ctx := context.TODO()
	registryInst, err := registry.InitRegistry(ctx, koalaConf.Register.RegisterName,
		registry.WithAddrs([]string{koalaConf.Register.RegisterAddr}),
		registry.WithHeartBeat(koalaConf.Register.HeartBeat),
		registry.WithRegistryPath(koalaConf.Register.RegisterPath),
		registry.WithTimeout(koalaConf.Register.Timeout),
	)
	if err != nil {
		logs.Error(ctx, "init register failed, err :%v", err)
		return
	}

	koalaServer.register = registryInst
	service := &registry.Service{
		Name: serviceName,
	}

	ip, err := util.GetLocalIP()
	if err != nil {
		logs.Error(ctx, "get local ip failed, err :%v", err)
		return
	}
	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   ip,
		Port: koalaConf.Port,
	})

	registryInst.Register(ctx, service)
	return
}

func Run() {
	if koalaConf.Prometheus.SwitchOn {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			addr := fmt.Sprintf("0.0.0.0:%d", koalaConf.Prometheus.Port)
			log.Fatal(http.ListenAndServe(addr, nil))
		}()
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", koalaConf.Port))
	if err != nil {
		log.Fatal("failed to listen :%v", err)
	}

	koalaServer.Serve(listen) // todo 为什么能直接调用serve
}

func GRPCServer() *grpc.Server {
	return koalaServer.Server // TODO 哪里来的Server?
}

func BuildServerMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware

	if koalaConf.Prometheus.SwitchOn {
		mids = append(mids, middleware.PrometheusServerMiddleware)
	}

	if koalaConf.Limit.SwitchOn {
		mids = append(mids, middleware.NewRateLimitMiddleware(koalaServer.limiter)) // TODO 这个limiter是怎么实现接口的
	}
	if len(koalaServer.userMiddleware) != 0 {
		mids = append(mids, koalaServer.userMiddleware...)
	}

	if len(mids) > 0 {
		// 把所有中间件组织成一个调用链
		m := middleware.Chain(mids[0], mids[1:]...)
		// 返回调用链的入口函数
		return m(handle)
	}

	m := middleware.Chain(mids[0])
	return m(handle)
}

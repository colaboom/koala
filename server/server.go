package server

import (
	"fmt"
	"github.com/koala/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type KoalaServer struct {
	*grpc.Server
	limiter *rate.Limiter

	userMiddleware []middleware.Middleware
}

var koalaServer = &KoalaServer{
	Server: grpc.NewServer(),
}

func Use(m ...middleware.Middleware) {
	koalaServer.userMiddleware = append(koalaServer.userMiddleware, m...)
}

func Init(serverName string) (err error) {
	err = InitConfig(serverName)
	if err != nil {
		return
	}

	if koalaConf.Limit.SwitchOn {
		koalaServer.limiter = rate.NewLimiter(rate.Limit(koalaConf.Limit.QPSLimit), koalaConf.Limit.QPSLimit)
	}

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

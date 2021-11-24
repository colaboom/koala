package server

import (
	"github.com/koala/middleware"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
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

func BuildServerMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware

	mids = append(mids, middleware.PrometheusServerMiddleware)

	if len(koalaServer.userMiddleware) != 0 {
		mids = append(mids, koalaServer.userMiddleware...)
	}

	if len(mids) > 0 {
		m := middleware.Chain(mids[0], mids[1:]...)
		return m(handle)
	}

	m := middleware.Chain(mids[0])
	return m(handle)
}

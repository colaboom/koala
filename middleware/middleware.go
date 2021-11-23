package middleware

import (
	"context"
)

type MiddlewareFunc func(ctx context.Context, req interface{}) (resp interface{}, err error)

type Middleware func(MiddlewareFunc) MiddlewareFunc

func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}

/*func BuildServerMiddleware(handle MiddlewareFunc) MiddlewareFunc {
	var mids []Middleware

	mids = append(mids, AccessLogMiddleware)
	if koalaConf.Prometheus.SwitchOn {
		mids = append(mids, PrometheusServerMiddleware)
	}

	if koalaConf.Limit.SwitchOn {
		mids = append(mids, NewRateLimitMiddleware(koalaServer.limiter))
	}

	if koalaConf.Trace.SwitchOn {
		mids = append(mids, TraceServerMiddleware)
	}

	if len(koalaServer.userMiddleware) != 0 {
		mids = append(mids, koalaServer.userMiddleware...)
	}

	m := Chain(PrepareMiddleware, mids...)
	return m(handle)
}*/
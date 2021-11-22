package middleware

import (
	"context"
	"github.com/koala/meta"
	"github.com/koala/middleware/prometheus"
	"time"
)

var (
	DefaultServerMetrics = prometheus.NewServerMetrics()
)

func PrometheusServerMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		serverMeta := meta.GetServerMeta(ctx)
		DefaultServerMetrics.IncrRequest(ctx, serverMeta.ServiceName, serverMeta.Method)

		startTime := time.Now()
		resp, err = next(ctx, req)

		DefaultServerMetrics.IncrCode(ctx, serverMeta.ServiceName, serverMeta.Method, err)
		DefaultServerMetrics.Latency(ctx, serverMeta.ServiceName, serverMeta.Method, time.Since(startTime).Nanoseconds()/1000)
		return
	}
}

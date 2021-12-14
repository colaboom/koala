package middleware

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/koala/meta"
)

func HystrixMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		rpcMeta := meta.GetRpcMeta(ctx)

		hystrixErr := hystrix.Do(rpcMeta.ServiceName, func() error {
			resp, err = next(ctx, req)
			return err
		}, nil)

		if hystrixErr != nil {
			return nil, hystrixErr
		}

		return resp, err
	}
}

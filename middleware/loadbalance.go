package middleware

import (
	"context"
	"github.com/koala/logs"
	"github.com/koala/errno"
	"github.com/koala/meta"
	"github.com/koala/loadbalance"
)

func NewLoadBalance(balancer loadbalance.LoadBalance) Middleware {
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			// 从ctx里获取rpcMeta
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) == 0 {
				err = errno.NotHaveInstance
				logs.Error(ctx, "not have instance")
				return
			}

			ctx = loadbalance.WithBalanceContext(ctx)
		}
	}
}
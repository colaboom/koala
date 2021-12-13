package middleware

import (
	"context"
	"github.com/koala/errno"
	"github.com/koala/loadbalance"
	"github.com/koala/logs"
	"github.com/koala/meta"
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
			for {
				rpcMeta.CurNode, err = balancer.Select(ctx, rpcMeta.AllNodes)
				if err != nil {
					return
				}
				logs.Info(ctx, "selected node:%#v", rpcMeta.CurNode)
				rpcMeta.HistoryNodes = append(rpcMeta.HistoryNodes, rpcMeta .CurNode)
				resp, err = next(ctx, req)
				if err != nil {
					// 连接错误的话，进行重试
					if errno.IsConnectError(err) {
						continue
					}
					return
				}
				break
			}
			return
		}
	}
}

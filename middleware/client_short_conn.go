package middleware

import (
	"context"
	"fmt"
	"github.com/koala/errno"
	"github.com/koala/logs"
	"github.com/koala/meta"
	"google.golang.org/grpc"
)

func ShortConnectMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		rpcMeta := meta.GetRpcMeta(ctx)
		if rpcMeta.CurNode == nil {
			err = errno.InvalidNode
			logs.Error(ctx, "invalid instance")
			return
		}

		addr := fmt.Sprintf("%s:%d", rpcMeta.CurNode.IP, rpcMeta.CurNode.Port)
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			logs.Error(ctx, "connect failed, err : %v", err)
			return nil, errno.ConnFailed
		}
		rpcMeta.Conn = conn
		defer conn.Close()
		resp, err = next(ctx, req)
		return
	}
}

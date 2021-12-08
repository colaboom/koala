package meta

import (
	"context"
	"github.com/koala/registry"
	"google.golang.org/grpc"
)

type RpcMeta struct {
	// 调用方名字
	Caller string
	// 提供方名字
	ServiceName string
	// 调用的方法
	Method string
	// 服务提供方的集群
	CallerCluster string
	// TraceID
	TraceID string
	// 环境
	Env string
	// 调用方IDC
	CallerIDC string
	// 服务提供方ICD
	ServiceIDC string
	// 当前节点
	CurNode *registry.Node
	// 历史选择节点
	HistoryNodes []*registry.Node
	// 服务提供方的节点列表
	AllNodes []*registry.Node
	// 当前请求使用的链接
	Conn *grpc.ClientConn
}

type rpcMetaContextKey struct{}

// TODO    server_meta和rpc_meta是用来干啥的
func GetRpcMeta(ctx context.Context) *RpcMeta {
	meta, ok := ctx.Value(rpcMetaContextKey{}).(*RpcMeta)
	if !ok {
		meta = &RpcMeta{}
	}

	return meta
}

func InitRpcMeta(ctx context.Context, service, method string) context.Context {
	meta := &RpcMeta{
		Method:      method,
		ServiceName: service,
	}

	return context.WithValue(ctx, rpcMetaContextKey{}, meta)
}

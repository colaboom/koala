package rpc

import (
	"context"
	"github.com/koala/loadbalance"
	"github.com/koala/logs"
	"github.com/koala/meta"
	"github.com/koala/middleware"
	"github.com/koala/registry"
	"sync"
	"time"
)

var initRegistryOnce sync.Once
var globalRegister registry.Registry

type KoalaClient struct {
	opts     *RpcOptions
	register registry.Registry
	balance  loadbalance.LoadBalance
}

func NewKoalaClient(serviceName string, optfunc ...RpcOptionFunc) *KoalaClient {
	client := &KoalaClient{
		opts: &RpcOptions{
			ConnTimeout:       DefaultConnTimeout,
			WriteTimeout:      DefaultWriteTimeout,
			ReadTimeout:       DefaultReadTimeout,
			ServiceName:       serviceName,
			RegisterName:      "etcd",
			RegisterAddr:      "127.0.0.1:2379",
			RegisterPath:      "/koala/service/",
			TraceReportAddr:   "http://60.205.218.189:9411/api/v1/spans",
			TraceSampleType:   "const",
			TraceSampleRate:   1,
			ClientServiceName: "default",
		},
		balance: loadbalance.NewRandomBalance(),
	}

	for _, opt := range optfunc {
		opt(client.opts)
	}

	initRegistryOnce.Do(func() {
		ctx := context.TODO()
		var err error
		globalRegister, err = registry.InitRegistry(ctx,
			client.opts.RegisterName,
			registry.WithAddrs([]string{client.opts.RegisterAddr}),
			registry.WithTimeout(time.Second),
			registry.WithRegistryPath(client.opts.RegisterPath),
			registry.WithHeartBeat(10),
		)
		if err != nil {
			logs.Error(ctx, "init registry failed, err:%v", err)
			return
		}
	})

	client.register = globalRegister

	return client
}

// TODO 为什么把serverMeta的serviceName作为caller
func (k *KoalaClient) getCaller(ctx context.Context) string {
	serverMeta := meta.GetServerMeta(ctx)
	logs.Info(ctx,"serverMeta.ServiceName : %v", serverMeta.ServiceName)
	if serverMeta == nil {
		return ""
	}
	return serverMeta.ServiceName
}

func (k *KoalaClient) buildMiddleware(handle middleware.MiddlewareFunc) middleware.MiddlewareFunc {
	var mids []middleware.Middleware
	mids = append(mids, middleware.NewDiscoveryMiddleware(k.register))
	mids = append(mids, middleware.NewLoadBalance(k.balance))
	mids = append(mids, middleware.ShortConnectMiddleware)

	m := middleware.Chain(mids[0], mids[1:]...)
	return m(handle)
}

func (k *KoalaClient) Call(ctx context.Context, method string, r interface{}, handle middleware.MiddlewareFunc) (resp interface{}, err error) {
	// 构建中间件
	caller := k.getCaller(ctx)
	ctx = meta.InitRpcMeta(ctx, k.opts.ServiceName, method, caller)
	middlewareFunc := k.buildMiddleware(handle)
	resp, err = middlewareFunc(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp, err

}

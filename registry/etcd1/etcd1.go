package etcd1

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/koala/registry"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxServiceNum          = 8
	MaxSyncServiceInterval = time.Second * 10
)

var (
	etcdRegistry = &EtcdRegistry{
		serviceCh:          make(chan *registry.Service, MaxServiceNum),
		registryServiceMap: make(map[string]*RegisterService, MaxServiceNum),
	}
)

type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service

	value              atomic.Value
	lock               sync.Mutex
	registryServiceMap map[string]*RegisterService
}

type RegisterService struct {
	id          clientv3.LeaseID
	service     *registry.Service
	registered  bool
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

func init() {
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	etcdRegistry.value.Store(allServiceInfo)
	registry.RegisterPlugin(etcdRegistry)
	go etcdRegistry.run()
}

func (er *EtcdRegistry) run() {
	for {
		select {
		case service := <- er.serviceCh:
			registerService, ok := er.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes{
					registerService.service.Nodes = append(registerService.service.Nodes, node)
				}
				registerService.registered = false // todo ????为啥是false
				break
			}
			registerService = &RegisterService{
				service: service,
			}
			er.registryServiceMap[service.Name] = registerService
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (er *EtcdRegistry) Name() string {
	return "etcd"
}

func (er *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {
	for _, opt := range opts {
		opt(er.options)
	}
	er.client, err = clientv3.New(clientv3.Config{
		Endpoints:   er.options.Addrs,
		DialTimeout: er.options.Timeout,
	})
	if err != nil {
		err = fmt.Errorf("init etcd failed , err : %v", err)
		return
	}
	return
}

func (er *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case er.serviceCh <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}
	return
}

func (er *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return
}

func (er *EtcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {
	return
}

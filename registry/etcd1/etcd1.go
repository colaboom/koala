package etcd1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/koala/registry"
	"path"
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
	ticker := time.NewTicker(MaxSyncServiceInterval)
	for {
		select {
		//将管道里面的service都同步到etcd
		case service := <-er.serviceCh:
			registerService, ok := er.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes {
					registerService.service.Nodes = append(registerService.service.Nodes, node)
				}
				// 初始化为false，只有注册到etcd中后才变为true
				registerService.registered = false
				break
			}
			registerService = &RegisterService{
				service: service,
			}
			er.registryServiceMap[service.Name] = registerService
		case <-ticker.C:
			// 定时从etcd中主动拉取最新的服务信息
			er.syncServiceFromEtcd()
		default:
			er.registerOrKeepAlive()
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (er *EtcdRegistry) registerOrKeepAlive() {
	for _, registryService := range er.registryServiceMap {
		if registryService.registered {
			er.keepAlive(registryService)
			continue
		}
		er.registerService(registryService)
	}
}

func (er *EtcdRegistry) keepAlive(registryService *RegisterService) {
	// 没有具体操作。从keepAliveCh管道拿数据，拿到了就说明正常已注册，没拿到就说明异常，讲注册状态改为false后重新注册
	select {
	case resp := <-registryService.keepAliveCh:
		if resp == nil {
			registryService.registered = false
			return
		}
		fmt.Printf("service:%s, ip:%s, port:%v\n", registryService.service.Name, registryService.service.Nodes[0].IP, registryService.service.Nodes[0].Port)
	}
	return
}

func (er *EtcdRegistry) registerService(registryService *RegisterService) (err error) {
	resp, err := er.client.Grant(context.TODO(), er.options.HeartBeat)
	if err != nil {
		return
	}

	registryService.id = resp.ID
	for _, node := range registryService.service.Nodes {
		tmp := &registry.Service{
			Name: registryService.service.Name,
			Nodes: []*registry.Node{
				node,
			},
		}
		key := er.serviceNodePath(tmp)
		data, err := json.Marshal(tmp)
		if err != nil {
			continue
		}
		_, err = er.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			continue
		}
		aliveCh, err := er.client.KeepAlive(context.TODO(), resp.ID)
		if err != nil {
			continue
		}

		registryService.keepAliveCh = aliveCh
		registryService.registered = true
	}
	return
}

func (er *EtcdRegistry) serviceNodePath(service *registry.Service) string {
	nodeIp := fmt.Sprintf("%s:%d", service.Nodes[0].IP, &service.Nodes[0].Port)
	return path.Join(er.options.RegistryPath, service.Name, nodeIp)
}

func (er *EtcdRegistry) Name() string {
	return "etcd"
}

func (er *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {
	er.options = &registry.Options{}
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
	// 从缓存中读取
	service, ok := er.getServiceFromCache(name)
	if ok {
		return
	}
	er.lock.Lock()
	defer er.lock.Unlock()
	service, ok = er.getServiceFromCache(name)
	if ok {
		return
	}

	// 如果缓存中没有，就从etcd中读取
	key := er.servicePath(name)
	resp, err := er.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return
	}
	service = &registry.Service{
		Name: name,
	}
	var tmpService registry.Service
	for _, kv := range resp.Kvs {
		value := kv.Value
		err = json.Unmarshal(value, &tmpService)
		if err != nil {
			return
		}
		for _, node := range tmpService.Nodes {
			service.Nodes = append(service.Nodes, node)
		}
	}

	// 更新缓存
	allServiceInfoNew := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}
	allServiceInfoOld := er.value.Load().(*AllServiceInfo)
	for key, val := range allServiceInfoOld.serviceMap {
		allServiceInfoNew.serviceMap[key] = val
	}
	allServiceInfoNew.serviceMap[name] = service
	er.value.Store(allServiceInfoNew)

	return
}

func (er *EtcdRegistry) getServiceFromCache(name string) (service *registry.Service, ok bool) {
	allServiceInfo := er.value.Load().(*AllServiceInfo)
	service, ok = allServiceInfo.serviceMap[name]

	return
}

// 定时从etcd中主动拉取最新的服务信息
func (er *EtcdRegistry) syncServiceFromEtcd() {
	var allServiceInfoNew = &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	}

	allServiceInfoOld := er.value.Load().(*AllServiceInfo)
	ctx := context.TODO()

	//对于缓存中的每一个服务，都需要从etcd中更新
	for _, service := range allServiceInfoOld.serviceMap {
		key := er.servicePath(service.Name)
		resp, err := er.client.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			allServiceInfoNew.serviceMap[service.Name] = service
			continue
		}

		serviceNew := &registry.Service{
			Name: service.Name,
		}
		for _, kv := range resp.Kvs {
			value := kv.Value
			var tmpService registry.Service
			err = json.Unmarshal(value, &tmpService)
			if err != nil {
				fmt.Printf("unmarshal failed, err:%v value:%s", err, string(value))
				return
			}
			for _, node := range tmpService.Nodes {
				serviceNew.Nodes = append(serviceNew.Nodes, node)
			}
		}
		allServiceInfoNew.serviceMap[service.Name] = serviceNew
	}

	er.value.Store(allServiceInfoNew)
}

func (er *EtcdRegistry) servicePath(name string) string {
	return path.Join(er.options.RegistryPath, name)
}

module github.com/koala

go 1.16

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.17+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/emicklei/proto v1.9.1
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/btree v1.0.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.11.3 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/urfave/cli v1.22.5
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.11.0 // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5
	golang.org/x/sys v0.0.0-20211117180635-dee7805ff2e1 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace google.golang.org/grpc v1.42.0 => google.golang.org/grpc v1.26.0

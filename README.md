# koala
基于grpc的golang微服务框架
- 代码自动生成
  - 生成客户端代码：koala.exe -c -f hello.proto
  - 生成服务端代码：koala.exe -s -f hello.proto
- 集成功能
  - 日志系统
  - 服务注册和服务发现
  - 负载均衡
  - 过载保护（熔断，限流）
  - 数据可视化监控（prometheus + grafana）

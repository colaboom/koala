package main

var config_template = `port: 8080
prometheus:
  switch_on: true
  port: 8081
service_name: {{.Package.Name}}
register:
  switch_on: true
  register_path: /koala/service/
  timeout: 1s
  heart_beat: 10
  register_name: etcd
  register_addr: 127.0.0.1:2379
log:
  level: debug
  path: ./logs/
limit:
  switch_on: true
  qps: 50000
trace:
  switch_on: true
  report_addr: http:xxx.com
  sample_type: const
  sample_rate: 1
`

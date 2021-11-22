package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type CounterLimit struct {
	counter      int64 // 计数器
	limit        int64 // 指定时间窗口内允许的最大请求数
	intervalNano int64 // 指定的时间窗口
	unixNano     int64 // unix时间戳，单位为纳秒
}

func NewCounterLimit(interval time.Duration, limit int64) *CounterLimit {
	return &CounterLimit{
		counter:      0,
		limit:        limit,
		intervalNano: int64(interval),
		unixNano:     time.Now().UnixNano(),
	}
}

func (c *CounterLimit) Allow() bool {
	now := time.Now().UnixNano()
	if now-c.intervalNano > c.unixNano {
		atomic.StoreInt64(&c.counter, 0)
		atomic.StoreInt64(&c.unixNano, now)
	}
	atomic.AddInt64(&c.counter, 1)
	return c.counter < c.limit // 判断是否允许请求进来
}

func main() {
	limit := NewCounterLimit(time.Second, 100)
	for i := 0; i < 200; i++ {
		ret := limit.Allow()
		fmt.Printf("这是第%d个请求，是否允许:%v\n", i, ret)
	}
}

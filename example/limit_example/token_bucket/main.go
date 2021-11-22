package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	limiter := rate.NewLimiter(50, 100)
	for i := 0; i < 1000; i++ {
		ret := limiter.Allow()
		time.Sleep(time.Millisecond * 5)
		if ret {
			fmt.Printf("第%d个请求，成功\n", i)
		} else {
			fmt.Printf("第%d个请求，失败\n", i)
		}
	}
}

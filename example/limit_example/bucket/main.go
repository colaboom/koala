package main

import (
	"fmt"
	"math"
	"time"
)

type BucketLimit struct {
	rate       float64 // 漏桶中水的流出速率
	bucketSize float64 // 漏桶最大的容量
	unixNano   int64   // unix时间戳
	curWater   float64 // 当前桶里面的水
}

func NewBucketLimit(rate float64, bucketSize float64) *BucketLimit {
	return &BucketLimit{
		rate:       rate,
		bucketSize: bucketSize,
		unixNano:   time.Now().UnixNano(),
		curWater:   0,
	}
}

func (b *BucketLimit) reflesh() {
	now := time.Now().UnixNano()
	diffSec := float64(now-b.unixNano) / 1000 / 1000 / 1000
	b.curWater = math.Max(0, b.curWater-diffSec*b.rate)
	b.unixNano = now
	return
}

func (b *BucketLimit) Allow() bool {
	b.reflesh()
	if b.curWater < b.bucketSize {
		b.curWater++
		return true
	}
	return false
}

func main() {
	limit := NewBucketLimit(50, 100)
	for i := 0; i < 1000; i++ {
		ret := limit.Allow()
		time.Sleep(time.Millisecond * 5)
		if ret {
			fmt.Printf("第%d个请求，成功\n", i)
		} else {
			fmt.Printf("第%d个请求，失败\n", i)
		}
	}
}

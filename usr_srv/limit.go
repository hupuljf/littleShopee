package main

import (
	"fmt"
	"sync"
	"time"
)

// 限流器结构体
type Limiter struct {
	rate        int        // 限流速率，单位：请求数/秒
	bucket      int        // 漏桶容量，单位：请求数
	water       int        // 漏桶中当前的水量，单位：请求数
	lastRequest time.Time  // 上次请求时间
	mutex       sync.Mutex // 互斥锁，用于并发安全
}

// 创建一个新的限流器
func NewLimiter(rate, bucket int) *Limiter {
	return &Limiter{
		rate:   rate,
		bucket: bucket,
		water:  0,
	}
}

// 判断是否允许通过
func (limiter *Limiter) Allow() bool {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	now := time.Now()
	// 计算当前时间与上次请求时间的时间间隔
	elapsed := now.Sub(limiter.lastRequest)

	// 计算漏桶中应该流出的水量
	drop := int(elapsed.Seconds() * float64(limiter.rate))

	// 更新上次请求时间
	limiter.lastRequest = now

	// 如果漏桶中的水量小于应该流出的水量，则将漏桶中的水量置为0
	if limiter.water < drop {
		limiter.water = 0
	} else {
		// 否则，漏桶中的水量减去应该流出的水量
		limiter.water -= drop
	}

	// 判断漏桶中的水量加上当前请求的水量是否超过漏桶的容量
	if limiter.water+1 > limiter.bucket {
		return false
	}

	// 更新漏桶中的水量
	limiter.water++
	return true
}

func main() {
	limiter := NewLimiter(2, 5) // 创建一个限流器，限制速率为2个请求/秒，漏桶容量为5个请求

	// 模拟请求
	for i := 1; i <= 10; i++ {
		time.Sleep(5000 * time.Millisecond) // 模拟请求间隔500ms
		if limiter.Allow() {
			fmt.Println("Request", i, "allowed")
		} else {
			fmt.Println("Request", i, "rejected")
		}
	}
}

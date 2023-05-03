package main

import (
	"fmt"
	tokenRateLimit "github.com/juju/ratelimit"
	"go.uber.org/ratelimit"
	"log"
	"time"
)

func main() {
	//leakyBucket()
	TokenBucket()
}

//漏桶策略
func leakyBucket() {
	//每一秒钟多少次请求
	rl := ratelimit.New(100)
	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
	//0 625ns
	//1 10ms
	//2 10ms
	//3 10ms
	//4 10ms
	//5 10ms
	//6 10ms
	//7 10ms
	//8 10ms
	//9 10ms

}
func TokenBucket() {
	//
	bucket := tokenRateLimit.NewBucket(time.Second, 10)
	for i := 0; i < 20; i++ {
		time.Sleep(time.Millisecond * 400)
		if bucket.TakeAvailable(1) < 1 {
			log.Println("take token failed")
		}
		log.Println("take token success")
	}

}

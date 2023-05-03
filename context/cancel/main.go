package main

import (
	"context"
	"log"
	"time"
)

func main() {
	//gen := func(ctx context.Context) <-chan int {
	//	dst := make(chan int)
	//	n := 1
	//	go func() {
	//		for {
	//			select {
	//			case <-ctx.Done():
	//				return // returning not to leak the goroutine
	//			case dst <- n:
	//				n++
	//			}
	//		}
	//	}()
	//	return dst
	//}
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel() // cancel when we are finished consuming integers
	//
	//for n := range gen(ctx) {
	//	fmt.Println(n)
	//	if n == 5 {
	//		break
	//	}
	//}
	//time.Sleep(time.Second * 5)
	//Output:
	//
	//1
	//2
	//3
	//4
	//5
	cancelParent()
}
func cancelParent() {
	//测试多个多个ctx嵌套的情况，第一层取消的情况
	ctx1, cancelFunc1 := context.WithCancel(context.Background())
	ctx2, cancelFunc2 := context.WithCancel(ctx1)
	defer cancelFunc2()
	go func() {
		time.Sleep(time.Second)
		cancelFunc1()
	}()
	<-ctx2.Done()
	log.Println(ctx2.Err())
	log.Println(ctx1.Err())
	log.Println("finish")
}

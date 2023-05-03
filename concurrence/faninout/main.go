package main

import (
	"fmt"
	"sync"
	"time"
)

//定义扇出（Fan-out）是一个术语，用于描述启动多个goroutines以处理来自管道的输入的过程，并且扇入（fan-in）是描述将多个结果组合到一个通道中的过程的术语。
//不依赖模块之前的计算结果。 应用场景
//运行需要很长时间。
func main() {
	tv := time.Now()

	firstResult := job1(10)

	// 拆分成三个job2，即3个goroutine （扇出）
	secondResult1 := job2(firstResult)
	secondResult2 := job2(firstResult)
	secondResult3 := job2(firstResult)

	// 合并结果(扇入）
	secondResult := merge(secondResult1, secondResult2, secondResult3)

	thirdResult := job3(secondResult)

	for v := range thirdResult {
		fmt.Println(v)
	}

	fmt.Println("all finish")
	fmt.Println("duration:", time.Since(tv).String())

}

func merge(inCh ...<-chan int) <-chan int {
	outCh := make(chan int, 2)

	var wg sync.WaitGroup
	for _, ch := range inCh {
		wg.Add(1)
		go func(wg *sync.WaitGroup, in <-chan int) {
			defer wg.Done()
			for val := range in {
				outCh <- val
			}
		}(&wg, ch)
	}

	// 重要注意，wg.Wait() 一定要在goroutine里运行，否则会引起deadlock
	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func job1(count int) <-chan int {
	outCh := make(chan int, 2)

	go func() {
		defer close(outCh)
		for i := 0; i < count; i++ {
			time.Sleep(time.Second)
			fmt.Println("job1 finish:", 1)
			outCh <- 1
		}
	}()

	return outCh
}

func job2(inCh <-chan int) <-chan int {
	outCh := make(chan int, 2)

	go func() {
		defer close(outCh)
		for val := range inCh {
			// 耗时2秒
			time.Sleep(time.Second * 2)
			val++
			fmt.Println("job2 finish:", val)
			outCh <- val
		}
	}()

	return outCh
}

func job3(inCh <-chan int) <-chan int {
	outCh := make(chan int, 2)

	go func() {
		defer close(outCh)
		for val := range inCh {
			val++
			fmt.Println("job3 finish:", val)
			outCh <- val
		}
	}()

	return outCh
}

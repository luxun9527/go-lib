package main

import (
	"fmt"
	"testing"
)

//运行当前 package 内的用例：go test example 或 go test
//go test -v benchmark_test.go -benchtime=5s  -bench=BenchmarkHello$ -count=3
//重要的参数
//go test -bench='Fib$' -benchtime=5s .  benchtime 指定运行的时间默认是运行1s 1, 2, 3, 5, 10, 20, 30, 50, 100
//-bench=BenchmarkHello$  指定匹配的测试用例
//--count=3表示测试三轮
// time.Sleep(3)
//b.ResetTimer() // 重置定时器 避免准备受到影响
//StopTimer & StartTime 还有一种情况，每次函数调用前后需要一些准备工作和清理工作，我们可以使用 StopTimer 暂停计时以及使用 StartTimer 开始计时
func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}
func BenchmarkHello1(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}
func RangeForMap1(m map[int]string) {
	for k, v := range m {
		_, _ = k, v
	}
}

const N = 1000

func initMap() map[int]string {
	m := make(map[int]string, N)
	for i := 0; i < N; i++ {
		m[i] = fmt.Sprint("www.flysnow.org", i)
	}
	return m
}

//go test -v benchmark_test.go -benchtime=5s  -bench=BenchmarkRangeForMap1$ -count=3
func BenchmarkRangeForMap1(b *testing.B) {
	m := initMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RangeForMap1(m)
	}
}

func initSlice() []string {
	s := make([]string, N)
	for i := 0; i < N; i++ {
		s[i] = "www.flysnow.org"
	}
	return s
}

//go test -v benchmark_test.go -benchtime=5s  -bench=BenchmarkForSlice$ -count=3
func BenchmarkForSlice(b *testing.B) {
	s := initSlice()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ForSlice(s)
	}
}

//go test -v benchmark_test.go -benchtime=5s  -bench=BenchmarkRangeForMap1$ -count=3
func BenchmarkRangeForSlice(b *testing.B) {
	s := initSlice()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RangeForSlice(s)
	}
}
func ForSlice(s []string) {
	len := len(s)
	for i := 0; i < len; i++ {
		_, _ = i, s[i]
	}
}

func RangeForSlice(s []string) {
	for i, v := range s {
		_, _ = i, v
	}
}

//go test -v benchmark_test.go -benchtime=5s  -bench=BenchmarkRangeArray$ -count=3
func BenchmarkRangeArray(b *testing.B) {
	m := [1000]int64{}
	for i := 0; i < 1000; i++ {
		m[i] = int64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(m); j++ {
			_ = m[j]
		}
	}
}

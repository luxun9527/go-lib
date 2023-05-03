package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

// fetchUrlDemo2 使用errgroup并发获取url内容
func fetchUrlDemo2() error {
	var g errgroup.Group
	var urls = []string{
		"http://pkg.go.dev",
		"http://www.liwenzhou.com",
		"http://www.yixieqitawangzhi.com",
	}
	for _, url := range urls {
		url := url // 注意此处声明新的变量
		// 启动一个goroutine去获取url内容
		g.Go(func() error {
			resp, err := http.Get(url)
			if err == nil {
				fmt.Printf("获取%s成功\n", url)
				resp.Body.Close()
			}
			return err // 返回错误
		})
	}
	if err := g.Wait(); err != nil {
		// 处理可能出现的错误
		fmt.Println(err)
		return err
	}
	fmt.Println("所有goroutine均成功")
	return nil
}
func main() {
	if err := fetchUrlDemo2(); err != nil {
		fmt.Println("err", err)
	}
}

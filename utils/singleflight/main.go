package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// 模拟缓存
var cache = make(map[string]string)
var cacheMutex sync.Mutex

// SingleFlight 组
var sfGroup singleflight.Group

// 获取数据的函数
func getData(key string) (string, error) {
	// 检查缓存
	cacheMutex.Lock()
	if val, ok := cache[key]; ok {
		cacheMutex.Unlock()
		return val, nil
	}
	cacheMutex.Unlock()

	// 使用 SingleFlight 避免重复请求
	result, err, _ := sfGroup.Do(key, func() (interface{}, error) {
		// 模拟数据库查询
		fmt.Printf("查询数据库: %s\n", key)
		time.Sleep(1 * time.Second) // 模拟查询耗时
		data := "Value for " + key

		// 写入缓存
		cacheMutex.Lock()
		cache[key] = data
		cacheMutex.Unlock()

		return data, nil
	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func main() {
	var wg sync.WaitGroup

	// 模拟多个请求同时查询同一个 key
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := getData("key1")
			if err != nil {
				fmt.Printf("错误: %v\n", err)
				return
			}
			fmt.Printf("获取到数据: %s\n", val)
		}()
	}

	wg.Wait()
}

package google

import (
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestErrgorup(t *testing.T) {
	var eg errgroup.Group
	var ret1, ret2, ret3 int
	now := time.Now()
	eg.Go(func() error {
		var err error
		ret1, err = api1()
		return err
	})
	eg.Go(func() error {
		var err error
		ret2, err = api2()
		return err
	})
	eg.Go(func() error {
		var err error
		ret3, err = api3()
		return err
	})
	err := eg.Wait()
	cost := time.Since(now)
	fmt.Printf("err:%v cost:%v ret1:%v ret2:%v ret3:%v\n", err, cost, ret1, ret2, ret3)
}
func api1() (int, error) {
	time.Sleep(time.Second)
	return 1, nil
}

func api2() (int, error) {
	time.Sleep(2 * time.Second)
	return 2, nil
}

func api3() (int, error) {
	time.Sleep(3 * time.Second)
	return 3, nil
}

var (
	g            singleflight.Group
	ErrCacheMiss = errors.New("cache miss")
)

func TestSingleflight(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)

	// 模拟10个并发
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			data, err := load("key")
			if err != nil {
				log.Print(err)
				return
			}
			log.Println(data)
		}()
	}
	wg.Wait()
}

// 获取数据
func load(key string) (string, error) {
	data, err := loadFromCache(key)
	if err != nil && err == ErrCacheMiss {
		// 利用 singleflight 来归并请求
		v, err, _ := g.Do(key, func() (interface{}, error) {
			data, err := loadFromDB(key)
			if err != nil {
				return nil, err
			}
			setCache(key, data)
			return data, nil
		})
		if err != nil {
			log.Println(err)
			return "", err
		}
		data = v.(string)
	}
	return data, nil
}

// getDataFromCache 模拟从cache中获取值 cache miss
func loadFromCache(key string) (string, error) {
	return "", ErrCacheMiss
}

// setCache 写入缓存
func setCache(key, data string) {}

// getDataFromDB 模拟从数据库中获取值
func loadFromDB(key string) (string, error) {
	fmt.Println("query db")
	unix := strconv.Itoa(int(time.Now().UnixNano()))
	return unix, nil
}

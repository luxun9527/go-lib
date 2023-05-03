package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//在Golang中，流水线由多个阶段组成，每个阶段之间通过channel连接，每个节点可以由多个同时运行的goroutine组成。
//应用场景：串行执行的任务，后面的任务依赖前面任务的结果，且任务都比较耗时。不需要返回返回值。比如多个任务都是io操作，频繁触发系统调用如果这时，在对接口性能有要求的场景可以使用流水线模型
type H map[string]interface{}

func main() {
	go func() {
		http.HandleFunc("/getUserNameInfo", func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Millisecond * 500)
			data := H{"name": "zhangsan", "ID": time.Now().UnixNano()}
			d, _ := json.Marshal(data)
			if _, err := writer.Write(d); err != nil {
				fmt.Printf("err = %v", err)
			}
		})
		// 直接使用 http 包的 ListenAndServe 函数监听服务
		http.ListenAndServe(":9999", nil)
	}()
	time.Sleep(time.Second * 3)
	pdc := fetchData()
	hdc := handleData(pdc)
	storeData(hdc)

}
func handleData(data <-chan []byte) <-chan []byte {
	out := make(chan []byte)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("err = %v", err)
			}
			close(out)
		}()
		for v := range data {
			time.Sleep(time.Millisecond * 400)
			var d H
			if err := json.Unmarshal(v, &d); err != nil {
				fmt.Println("err", err)
				continue
			}
			fmt.Printf("data = %v \n", d)
			out <- v
		}

	}()
	return out
}
func fetchData() <-chan []byte {
	out := make(chan []byte)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("err = %v", err)
			}
			close(out)
		}()
		for {
			resp, err := http.Get("http://localhost:9999/getUserNameInfo")
			if err != nil {
				return
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				fmt.Println("err", err)
				return
			}
			resp.Body.Close()
			out <- data
		}
	}()
	return out
}
func storeData(data <-chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("err = %v", err)
		}
	}()
	for v := range data {
		time.Sleep(time.Millisecond * 400)
		file, err := os.OpenFile("data.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		v = append(v, []byte{'\n'}...)
		if _, err := file.Write(v); err != nil {
			fmt.Println("err", err)
			continue
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	data := map[string]interface{}{"name": "zhangsan", "age": 23}
	d, _ := json.Marshal(data)
	var output string
	for i := 0; i < 10; i++ {
		output += string(d)
	}
	for i := 0; true; i++ {
		time.Sleep(time.Second)
		fmt.Println(output, i)
	}
}

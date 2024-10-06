package main

import (
	"context"
	"log"
	"os/exec"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", `echo test`)
	//标准输出和标准错误输出都连接到同一个管道,获取标准输出和标准错误输出
	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("exec shell failed %v", err)
	}
	log.Println(string(data))

	cmd = exec.CommandContext(ctx, "sh", "-c", `cat a.txt`)
	//如果退出码不是0 ，则会有error
	data, err = cmd.CombinedOutput()
	if err != nil {
		log.Panicf("exec shell failed %v", err)
	}
}

package main

import (
	"bytes"
	"context"
	"log"
	"os/exec"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", `cat a.txt`)
	//cmd.CombinedOutput()
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	data := b.Bytes()
	if err != nil {
		log.Panicf("exec shell failed err %v detail %v", err, string(data))
	}
	//我们可以指定一个命令的标准输出和标准错误输出。
	//如果这个命令的退出码不是0，那么Run函数会返回一个错误。并且详细的错误信息会从stderr标准错误输出中获取。
	log.Printf("data:%v", string(data))

}

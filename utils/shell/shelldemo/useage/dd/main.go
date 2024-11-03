package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// 创建一个 bytes.Buffer 来保存输出
	var outputBuffer bytes.Buffer

	// 创建 exec.Command，使用 sh -c 执行命令
	cmd := exec.Command("sh", "-c", "dd if=/dev/stdin of=output_file.txt")

	// 获取标准输出和标准错误输出的管道
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer

	// 获取标准输入的管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Error creating stdin pipe: %v", err)
		return
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v", err)
		return
	}

	// 向 stdin 发送数据
	_, err = stdin.Write([]byte("Hello, world!\nThis is a test.\n"))
	if err != nil {
		log.Printf("Error writing to stdin: %v", err)
		return
	}

	// 关闭 stdin，模拟 Ctrl+D
	if err := stdin.Close(); err != nil {
		log.Printf("Error closing stdin: %v", err)
		return
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting for command: %v data %v", err, outputBuffer.String())
		return
	}

	// 打印输出和错误
	log.Println(outputBuffer.String())
}

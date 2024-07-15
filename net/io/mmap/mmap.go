package main

import (
	"fmt"
	"golang.org/x/exp/mmap"
	"log"
)

func main() {
	// 打开文件进行内存映射
	reader, err := mmap.Open("example.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer reader.Close()

	// 获取文件大小
	size := reader.Len()
	fmt.Printf("File size: %d bytes\n", size)

	// 读取文件内容
	data := make([]byte, size)
	_, err = reader.ReadAt(data, 0)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// 打印文件内容
	fmt.Println("File content:")
	fmt.Println(string(data))
}

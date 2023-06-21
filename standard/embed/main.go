package main

import (
	"embed"
	"fmt"
)

//go:embed test.txt
var embededFiles embed.FS

func main() {
	data, err := embededFiles.ReadFile("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data)) // hello world!
}

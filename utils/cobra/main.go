package main

import (
	"go-lib/utils/cobra/cmd"
)

func main() {
	cmd.Execute()
}

/*
go build o hugo.exe
核心概念

一个条命令有三个部分
命令自身和子命令
参数
flag
示例
hugo  times -t=10 34
hugo  times 为命令
-t=10 为flag
34 为参数
*/

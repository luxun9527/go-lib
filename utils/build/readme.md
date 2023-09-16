# 1、go的编译参数

##  -ldflags 参数

```go
go build -ldflags "-s -w" x.go
```

这两个参数的作用是让程序体积变小

https://www.cnblogs.com/bing-l/p/4072710.html?utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com

**-s** 去掉符号表，然后 panic 的时候 stack trace 就没有任何文件名/行号信息了，这个等价于普通C/C++程序被strip的效果。**不推荐使用**。存疑

**在实际的测试中没有发现panic实际的堆栈信息有什么不同**

```go

package main

import (
	"errors"
)
func main(){
	p1()
}
func p1(){
	p()
}
func p(){
	panic(errors.New("unknown error"))
}
```

![image-20230916111224914](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20230916111224914.png)

https://www.jianshu.com/p/095f921ca243

1、符号表存储了全局变量和函数变量名对应的地址。

2、符号表只在编译时期使用，删除符号表无法调试。

如果不调试，且对编译程序大小有严格要求，在生产环境加上这个参数没什么问题。

**-w** 去掉 DWARF 调试信息，得到的程序就不能用 gdb 调试了，如果不打算用 gdb 调试，基本没啥损失。



## -gcflags

```bash
go build -gcflags "-N -l"
```

-N代表禁止优化，不要在生产环境上开启，此处仅为演示使用

-l参数代表禁止内联，也建议不要在生产环境上开启，此处仅为演示使用

用在远程调试的情况

go build -gcflags "all=-N -l" github.com/app/demo
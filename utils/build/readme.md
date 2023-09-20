# 1、go的编译参数

##  -ldflags 参数

### -s -w

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
//打印的堆栈
==========================================================
$ go build   main.go && ./main.exe
panic: unknown error

goroutine 1 [running]:
main.p(...)
        E:/GoCode/go-lib/utils/build/main.go:18
main.p1(...)
        E:/GoCode/go-lib/utils/build/main.go:15
main.main()
        E:/GoCode/go-lib/utils/build/main.go:11 +0x49

================================================================

$ go build -ldflags "-s -w" main.go && ./main.exe
panic: unknown error                                 
                                                     
goroutine 1 [running]:                               
main.p(...)                                          
        E:/GoCode/go-lib/utils/build/main.go:17      
main.p1(...)                                         
        E:/GoCode/go-lib/utils/build/main.go:14      
main.main()                                          
        E:/GoCode/go-lib/utils/build/main.go:11 +0x49

```







https://www.jianshu.com/p/095f921ca243

1、符号表存储了全局变量和函数变量名对应的地址。

2、符号表只在编译时期使用，删除符号表无法调试。

如果不调试，且对编译程序大小有严格要求，在生产环境加上这个参数没什么问题。

**-w** 去掉 DWARF 调试信息，得到的程序就不能用 gdb 调试了，如果不打算用 gdb 调试，基本没啥损失。



### -x注入参数

```go
package inject

import (
	"flag"
	"fmt"
	"os"
)

var (
	builtAt   string
	buildUser string
	builtOn   string
	goVersion string
	gitAuthor string
	gitCommit string
	gitTag    string
)

func init(){
	flag.Bool("version",false,"打印版本信息")
	Register("version",PrintVersionInfo)
}
func PrintVersionInfo(val string) {
	fmt.Printf("%-20s %s\n", "builtAt", builtAt)
	fmt.Printf("%-20s %s\n", "builtOn", builtOn)
	fmt.Printf("%-20s %s\n", "buildUser", buildUser)
	fmt.Printf("%-20s %s\n", "goVersion", goVersion)
	fmt.Printf("%-20s %s\n", "gitAuthor", gitAuthor)
	fmt.Printf("%-20s %s\n", "gitCommit", gitCommit)
	fmt.Printf("%-20s %s\n", "gitTag", gitTag)
	os.Exit(1)
}


```



```go
package inject

import (
	"flag"
	"os"
)

//面对的情况，当设置了某个flag，就执行某个函数。
var	defaultFlag = &Flag{
	FlagSet: flag.CommandLine,
	m:       make(map[string]func(val string), 5),
}
type Flag struct {
	*flag.FlagSet
	m map[string]func(val string)
}

func Register(name string, handler func(val string)) {
	defaultFlag.register(name, handler)
}
func Parse() {
	defaultFlag.parseFlag()
}

func (f *Flag) register(name string, handler func(val string)) {
	f.m[name] = handler
}
func (f *Flag) parseFlag() {
	f.Parse(os.Args[1:])
	f.Visit(func(fl *flag.Flag) {
		f1, ok := f.m[fl.Name]
		if !ok {
			return
		}
		f1(fl.Value.String())
	})
}

```



## -gcflags

```bash
go build -gcflags "-N -l"
```

-N代表禁止优化，不要在生产环境上开启，此处仅为演示使用

-l参数代表禁止内联，也建议不要在生产环境上开启，此处仅为演示使用

用在远程调试的情况

go build -gcflags "all=-N -l" github.com/app/demo
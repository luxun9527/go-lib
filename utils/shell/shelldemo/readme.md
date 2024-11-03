# 核心概念

## 标准输入和标准输出

```plain
var (
    Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
    Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
    Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)
```

## 退出码

https://cloud.tencent.com/developer/article/1464331

shell命令的退出状态码都有特殊的意义，用来显示命令退出时的状态，更多地给外部使用．shell退出状态码是一个0~255之间的整数值．通常成功返回0，失败返回非0(错误码)．

1. 一般的退出状态码含义:

| **状态码** | **描述**                                                |
| ---------- | ------------------------------------------------------- |
| 0          | 命令成功结束                                            |
| 1          | 通用未知错误                                            |
| 2          | 误用Shell命令                                           |
| 126        | 命令不可执行                                            |
| 127        | 没找到命令                                              |
| 128        | 无效退出参数                                            |
| 128+n      | Linux信号n的致命错误 例: kill -9 ppid 出错返回128+9=137 |
| 130        | 命令通过Ctrl+C控制码越界                                |
| 255        | 退出码越界                                              |

echo $?

获取上一条命令执行的退出码

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1730045134714-d0e33fad-dc30-4b7f-8fe9-6686e0c96769.png)



# go执行shell

## 基础用法

```go
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

    //2024/10/28 00:04:08 exec shell failed err exit status 1 detail cat: a.txt: No such file or directory
}
```

## 写入到标准输入

```go
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
```
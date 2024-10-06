# 核心概念

## 标准输入和标准输出

```go
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

| 状态码 | 描述                                                    |
| :----- | :------------------------------------------------------ |
| 0      | 命令成功结束                                            |
| 1      | 通用未知错误                                            |
| 2      | 误用Shell命令                                           |
| 126    | 命令不可执行                                            |
| 127    | 没找到命令                                              |
| 128    | 无效退出参数                                            |
| 128+n  | Linux信号n的致命错误 例: kill -9 ppid 出错返回128+9=137 |
| 130    | 命令通过Ctrl+C控制码越界                                |
| 255    | 退出码越界                                              |

```javascript
echo $?
```

获取上一条命令执行的退出码
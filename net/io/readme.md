# go io操作

**1、相关概念**

**2、go相关的sdk**

## 理论知识

https://wxler.github.io/2021/02/19/134758/

### 相关概念

io操作在日常开发中接触的基本上都是磁盘io和网络io

input    就是将磁盘或网卡中的数据读取到我们用户空间的内存中。

output 将用户空间内存中的数据写到磁盘或网卡中。

具体流程是 读：磁盘/网卡-->内核缓冲区-->用户空间，写则相反。

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1720876642675-7ede38d7-d303-49b7-bb0c-b847de244845.png)

### 拓展知识dma

[https://wxler.github.io/2021/02/19/134758/#4-io%E4%B8%AD%E6%96%AD%E4%B8%8Edma](https://wxler.github.io/2021/02/19/134758/#4-io中断与dma)

从磁盘/网卡拷贝数据到内核缓冲区的工作由cpu负责，转为磁盘/网卡自己来做，提高cpu的效率。



### io类型

概念：文件描述符。linux系统中一切都是文件， 文件描述符（File Descriptor, FD）是操作系统内核为每个打开的文件分配的一个非负整数，用于唯一标识该文件 。可以抽象理解为操作系统读取文件的接口。

**如何让操作系统知道该文件描述符已经准备好了被读取和被写入？**

1、BIO （Blocking I/O），称之为同步阻塞I/O ，每个文件描述符启动一个线程，检查该文件描述符是否就绪，如果没就是就阻塞等待该文件描述符就绪。

2、NIO （Non-blocking IO），称之为非阻塞IO，只启动一个线程，不断检查所有文件描述符是否就绪，如果没有就绪就立刻返回，而检查的过程会进行系统调用消耗资源

**3、IO多路复用。**多路复用IO把轮询多个文件描述符放在内核空间里执行，即让内核负责监听所有的文件描述符（这样就不会有用户态和系统态的切换），当有文件描述符就绪，就将这些文件描述符返回给用户进程。


 ![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1720883187077-ccc521d3-88c1-4bdc-95ef-3f643ebbc4ce.png)

多路复用模式包含三种，即select、poll和epoll，**这几种模式主要区别在于对就绪的文件描述符的方式**。



**select** 最多能将1024个文件描述符加入到监听的集合中，当有文件描述符就绪后立刻返回。select只知道有文件描述符就绪了，但不知道是哪个，任然需要变量整个集合找出就绪的。

**poll**  poll和select类似，但是poll使用链表存储文件描述符，没有1024的限制。

**epoll** 使用红黑树，存储监听的文件描述符。当文件描述符就绪后使用链表存储就绪的文件描述符。不用像select和poll遍历所有监听的文件描述符找出就绪的, epoll的工作模式 水平触发（LT模式）,边缘触发（ET模式） 

LT如何文件描述符就绪，你不处理则会一直通知，ET只会通知一次。











**4、****AIO （ Asynchronous I/O）**：异步非阻塞I/O模型。传输过程如下：


 ![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1720884448690-f2015723-96c2-4f27-996d-3e4e80cd7564.png)

可以看到，异步非阻塞I/O在判断数据有没有准备好（即文件描述符是否就绪）和真正读数据两个阶段都是非阻塞的。AIO在第一次执行系统调用后，会注册一个回调函数，内核在检测到某文件描述符是否就绪，调用该回调函数执行真正的读操作，将数据从内核空间拷贝到用户空间，然后返回给用户使用。在整个过程，用户进程都是非阻塞状态，可以做其它的事情。类似js的回调函数。

### 零拷贝

将一个文件通过网络发送的流程：

磁盘数据-->系统内核缓冲区-->用户空间-->系统内核缓冲区-->网卡.涉及多次将数据四次缓冲区的拷贝，并经历了四次内核态和用户态的切换。

下图是不使用zero copy的网络IO传输过程：

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1720886533445-d2657a64-22d8-4122-84dd-b24da26d8d57.png)

零拷贝则不用将数据拷贝到用户空间。直接在两个内核缓冲区之间拷贝。

```
ssize_t sendfile(int out_fd, int in_fd, off_t *offset, size_t count)
```



# go相关sdk

### io多路复用

go 主要是使用

"golang.org/x/sys/unix"， "syscall"这两个包



**使用epoll创建一个tcp服务器**

```go
//go:build linux
// +build linux

package main

import (
    "errors"
    "log"
    "syscall"
)

func main() {
    // 创建 socket
    fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
    if err != nil {
        log.Fatalf("Error creating socket: %s", err)
    }
    defer syscall.Close(fd)

    //// 设置 socket 为非阻塞模式
    if err = syscall.SetNonblock(fd, true); err != nil {
        log.Fatalf("Error setting non-blocking mode: %s", err)
    }

    // 绑定地址和端口
    addr := syscall.SockaddrInet4{Port: 8080}
    copy(addr.Addr[:], []byte{0, 0, 0, 0}) // 绑定到所有 IP 地址
    if err = syscall.Bind(fd, &addr); err != nil {
        log.Fatalf("Error binding socket: %s", err)
    }

    // 开始监听
    if err = syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
        log.Fatalf("Error listening on socket: %s", err)
    }

    // 创建 epoll 实例
    epfd, err := syscall.EpollCreate1(0)
    if err != nil {
        log.Fatalf("Error creating epoll instance: %s", err)
    }
    defer syscall.Close(epfd)

    // 将监听 socket 添加到 epoll 实例
    if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
        Events: syscall.EPOLLIN,
        Fd:     int32(fd),
    }); err != nil {
        log.Fatalf("Error adding listener to epoll: %s", err)
    }

    // 事件循环
    events := make([]syscall.EpollEvent, 10)
    for {
        n, err := syscall.EpollWait(epfd, events, -1)
        if err != nil {
            if errors.Is(err, syscall.EINTR) {
                log.Printf("Epoll wait was interrupted, retrying...")
                continue
            }
            log.Fatalf("Error waiting for epoll events: %s", err)
        }
        log.Printf("Received %d events", n)
        for i := 0; i < n; i++ {
            if events[i].Fd == int32(fd) {
                // 接受新连接
                nfd, _, err := syscall.Accept(fd)
                if err != nil {
                    if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) {
                        log.Printf("No connections are available to accept")
                        continue
                    }
                    log.Printf("Error accepting connection: %s", err)
                    continue
                }

                // 设置 socket 为非阻塞模式
                if err = syscall.SetNonblock(nfd, true); err != nil {
                    log.Fatalf("Error setting non-blocking mode: %s", err)
                }
                // 将新连接添加到 epoll 实例
                if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, nfd, &syscall.EpollEvent{
                    Events: syscall.EPOLLIN | -syscall.EPOLLET,
                    Fd:     int32(nfd),
                }); err != nil {
                    log.Printf("Error adding connection to epoll: %s", err)
                    syscall.Close(nfd)
                    continue
                }

                log.Printf("Accepted new connection")
            } else {
                log.Printf("Ready for connection %+v", events[i])
                // 处理已就绪的连接
                cfd := int(events[i].Fd)
                go handleConnection(cfd, epfd)
            }
        }
    }
}

func handleConnection(cfd int, epfd int) {
    // defer syscall.Close(cfd)
    buf := make([]byte, 1024)

    for {
        n, err := syscall.Read(cfd, buf)
        if err != nil {
            log.Printf("Error reading from connection: %s", err)
            if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EBADF) {
                // 非阻塞模式下无数据可读
                break
            }
            return
        }
        if n == 0 {
            // 连接已关闭
            log.Printf("Connection closed")
            return
        }
        log.Printf("Received data: %s", string(buf[:n]))
        // 回显数据
        if _, err = syscall.Write(cfd, buf[:n]); err != nil {
            log.Printf("Error writing to connection: %s", err)
            return
        }
        return
    }
}
```

### 零拷贝

#### sendfile

sendfile 实现大文件不经过用户空间，直接传输。

```go
//go:build linux
// +build linux

package main

import (
    "fmt"
    "golang.org/x/sys/unix"
    "log"
    "net/http"
    "os"
)

// UploadHandler handles file upload and saves it using sendfile
func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // 10mb以上存临时文件，否则存储内存。
    err := r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        log.Printf("Error parsing form: %v", err)
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }

    // Get the file from form data
    file, handler, err := r.FormFile("file")
    if err != nil {
        log.Printf("Error retrieving file: %v", err)
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create a destination file
    dstPath := handler.Filename
    dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
    if err != nil {
        log.Printf("Error creating destination file: %v", err)
        http.Error(w, "Error creating destination file", http.StatusInternalServerError)
        return
    }
    defer dstFile.Close()

    // Get file descriptors
    srcFd := int(file.(*os.File).Fd())
    dstFd := int(dstFile.Fd())

    // Get the file size
    fi, err := file.(*os.File).Stat()
    if err != nil {
        log.Printf("Error getting file info: %v", err)
        http.Error(w, "Error getting file info", http.StatusInternalServerError)
        return
    }
    fileSize := fi.Size()

    // Use sendfile to transfer data
    _, err = unix.Sendfile(dstFd, srcFd, nil, int(fileSize))
    if err != nil {
        log.Printf("Error using sendfile: %v", err)
        http.Error(w, "Error using sendfile", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func main() {
    http.HandleFunc("/upload", UploadHandler)
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}
```



### mmap

`mmap` 是内存映射文件（Memory Mapped File）的缩写，是一种将文件的内容映射到进程的虚拟地址空间中的技术。通过 `mmap`，应用程序可以将文件或设备（如物理内存）的内容与内存地址关联，从而可以像访问内存一样访问文件内容。这种技术可以提高文件读写的效率，特别是对于大文件的操作，因为它避免了将文件内容拷贝到用户空间的过程。  

```go
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
```
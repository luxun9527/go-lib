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
	//	defer syscall.Close(cfd)
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

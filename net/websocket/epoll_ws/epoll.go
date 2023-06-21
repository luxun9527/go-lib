package main

import (
	"errors"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"reflect"
	"syscall"
)

var Epoller *Epoll

type Epoll struct {
	fd int64
	f  func()
}

func Start() {
	var err error
	Epoller, err = MkEpoll()
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			events := make([]unix.EpollEvent, 100)
			n, err := unix.EpollWait(int(Epoller.fd), events, -1)
			if errors.Is(err, syscall.EINTR) {
				continue
			}
			if errors.Is(err, syscall.EAGAIN) {
				break
			}

			for i := 0; i < n; i++ {
				if events[i].Events == syscall.EPOLLIN|syscall.EPOLLRDHUP {
					//Epoller.Remove(fd)
					continue
				}
				fd := int(events[i].Fd)
				CM.Notify(fd)
			}

		}

	}()
}
func MkEpoll() (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &Epoll{
		fd: int64(fd),
	}, nil
}

func (e *Epoll) Add(conn net.Conn) error {
	// Extract file descriptor associated with the connection
	fd := WebsocketFD(conn)
	//使用et模式
	if err := unix.EpollCtl(int(e.fd), syscall.EPOLL_CTL_ADD, int(fd), &unix.EpollEvent{Events: unix.EPOLLERR | unix.EPOLLET | unix.EPOLLRDHUP | unix.EPOLLPRI | unix.EPOLLIN, Fd: int32(fd)}); err != nil {
		return err
	}

	return nil
}
func WebsocketFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}

func (e *Epoll) Remove(conn net.Conn) error {
	fd := WebsocketFD(conn)
	if err := unix.EpollCtl(int(e.fd), syscall.EPOLL_CTL_DEL, int(fd), nil); err != nil {
		return err
	}
	return nil
}

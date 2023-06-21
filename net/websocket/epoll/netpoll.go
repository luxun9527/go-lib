//go:build linux
// +build linux

package main

import (
	"github.com/gobwas/ws"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func main() {

	listener, err := NewListener("127.0.0.1", 8989)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := listener.Accept()
	u := ws.Upgrader{
		OnHeader: func(key, value []byte) (err error) {
			log.Printf("non-websocket header: %q=%q", key, value)
			return
		},
	}
	if _, err := u.Upgrade(conn); err != nil {
		log.Fatalln(err)
	}
}

type Listener struct {
	fd   int
	addr net.Addr         // listener's local addr
	ln   *net.TCPListener // tcp|unix listener
	file *os.File
}

func (listener *Listener) Accept() (*NetFD, error) {

	connFd, addr, err := syscall.Accept(listener.fd)
	if err != nil {
		return nil, err
	}
	log.Println(addr)
	return &NetFD{
		fd: connFd,
		ip: "",
	}, nil

}
func NewListener(ip string, port int) (*Listener, error) {

	tcpListener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.ParseIP(ip).To4(),
		Port: port,
		Zone: "",
	})
	if err != nil {
		return nil, err
	}
	fileFd, err := tcpListener.File()
	if err != nil {
		return nil, err
	}
	return &Listener{
		fd:   int(fileFd.Fd()),
		addr: tcpListener.Addr(),
		ln:   tcpListener,
		file: fileFd,
	}, nil

}

type NetFD struct {
	fd int
	ip string
}

// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
// Read implements Conn.
func (netFd *NetFD) Read(b []byte) (n int, err error) {
	n, err = syscall.Read(netFd.fd, b)
	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EINTR {
			return 0, nil
		}
	}
	return n, err
}

// Write implements Conn.
func (netFd *NetFD) Write(b []byte) (n int, err error) {
	n, err = syscall.Write(netFd.fd, b)
	if err != nil {
		if err == syscall.EAGAIN {
			return 0, nil
		}
	}
	return n, err
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (netFd *NetFD) Close() error {
	return nil
}

// LocalAddr returns the local network address, if known.
func (netFd *NetFD) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr returns the remote network address, if known.
func (netFd *NetFD) RemoteAddr() net.Addr {
	return nil
}

func (netFd *NetFD) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (netFd *NetFD) SetReadDeadline(t time.Time) error {
	return nil
}

func (netFd *NetFD) SetWriteDeadline(t time.Time) error {
	return nil
}

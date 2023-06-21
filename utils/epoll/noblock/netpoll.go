//go:build linux
// +build linux

package main

import (
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

func main() {

	listener, err := NewListener("0.0.0.0", 8989)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := listener.Accept()
	//u := ws.Upgrader{
	//	OnHeader: func(key, value []byte) (err error) {
	//		log.Printf("non-websocket header: %q=%q", key, value)
	//		return
	//	},
	//}
	//if _, err := u.Upgrade(conn); err != nil {
	//	log.Fatalln(err)
	//}
	for {
		//	ws.ReadFrame()
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(buf[:n]))
	}
}

type Listener struct {
	fd   int
	addr net.Addr         // listener's local addr
	ln   *net.TCPListener // tcp|unix listener
	file *os.File
}

func (listener *Listener) Accept() (*NetFD, error) {

	connFd, sa, err := syscall.Accept(listener.fd)
	if err != nil {
		return nil, err
	}
	addr := sockaddrToAddr(sa)
	log.Println(addr.String())
	return &NetFD{
		fd: connFd,
		ip: "",
	}, nil

}
func sockaddrToAddr(sa syscall.Sockaddr) net.Addr {
	var a net.Addr
	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		a = &net.TCPAddr{
			IP:   sa.Addr[0:],
			Port: sa.Port,
		}
	case *syscall.SockaddrInet6:
		var zone string
		if sa.ZoneId != 0 {
			if ifi, err := net.InterfaceByIndex(int(sa.ZoneId)); err == nil {
				zone = ifi.Name
			}
		}
		// if zone == "" && sa.ZoneId != 0 {
		// }
		a = &net.TCPAddr{
			IP:   sa.Addr[0:],
			Port: sa.Port,
			Zone: zone,
		}
	case *syscall.SockaddrUnix:
		a = &net.UnixAddr{Net: "unix", Name: sa.Name}
	}
	return a
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

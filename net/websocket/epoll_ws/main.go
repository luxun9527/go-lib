package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	gws "github.com/gobwas/ws"
	"github.com/gobwas/ws/wsflate"
	"io"
	"log"
	"net"
	"sync"
	"syscall"
	"time"
)

func main() {
	engine := gin.New()
	NewManger()
	Start()
	NewBucket()
	engine.GET("/ws", Connect)
	engine.Run(":8989")
}

type Connection struct {
	buf  *bytes.Buffer
	conn net.Conn
	id   int
}

func NewConnection(conn net.Conn) *Connection {
	fd := WebsocketFD(conn)
	connection := &Connection{
		buf:  bytes.NewBuffer(make([]byte, 0, 100)),
		conn: conn,
		id:   fd,
	}
	return connection
}

type Bucket struct {
	c           []chan int
	Connections map[int]*Connection
	lock        sync.RWMutex
}

func (b *Bucket) AddConnection(connection *Connection) {
	b.lock.Lock()
	b.Connections[connection.id] = connection
	b.lock.Unlock()
}

type Manger struct {
	buckets []*Bucket
}

var CM *Manger

func (m *Manger) getBucket(connID int) *Bucket {
	i := connID % 55
	return m.buckets[i]
}
func (m *Manger) AddConnection(connection *Connection) {
	bucket := m.getBucket(connection.id)
	bucket.AddConnection(connection)
}
func (m *Manger) Notify(fd int) {
	bucket := m.getBucket(fd)
	i := fd % len(bucket.c)
	bucket.c[i] <- fd
}

func NewManger() {
	buckets := make([]*Bucket, 55)
	for i := 0; i < len(buckets); i++ {
		buckets[i] = NewBucket()
	}
	manger := &Manger{}
	manger.buckets = buckets
	CM = manger
}

func NewBucket() *Bucket {
	b := &Bucket{}
	b.Connections = make(map[int]*Connection, 50)
	cs := make([]chan int, 20)
	for i := 0; i < len(cs); i++ {
		cs[i] = make(chan int)
		go b.handleMessage(cs[i])
	}
	b.c = cs
	return b
}
func (b *Bucket) handleMessage(fds chan int) {
	for fd := range fds {
		b.lock.Lock()
		conn, ok := b.Connections[fd]
		b.lock.Unlock()
		if !ok {
			continue
		}
		conn.read(fd)
	}
}
func (c *Connection) read(fd int) {
	for {
		buf := make([]byte, 100)
		n, err := syscall.Read(fd, buf)
		if err != nil {
			//没有读到不阻塞
			if err == syscall.EAGAIN {
				break
			}
		}
		c.buf.Write(buf[:n])
		//循环读
		for {
			frame, err := gws.ReadFrame(c.buf)
			if err != nil {
				if err == io.ErrUnexpectedEOF {
					//长度没达到就写回去
					c.buf.Write(buf[:n])
					break
				}
				if err == io.EOF {
					break
				}

			}
			frame = gws.UnmaskFrameInPlace(frame)
			compressed, err := wsflate.IsCompressed(frame.Header)
			if err != nil {
				return
			}
			if compressed {
				frame, err = wsflate.DecompressFrame(frame)
				if err != nil {
					return
				}
			} else {
				log.Println(string(frame.Payload))
			}
		}

	}
}

func Connect(c *gin.Context) {
	var httpUpgrade gws.HTTPUpgrader
	conn, _, _, err := httpUpgrade.Upgrade(c.Request, c.Writer)
	if err != nil {
		return
	}
	if err := Epoller.Add(conn); err != nil {
		log.Println("err", err)
		return
	}
	connection := NewConnection(conn)
	CM.AddConnection(connection)
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
}

package object

import (
	"errors"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

type myint int32

//https://cyent.github.io/golang/basic/type_custom/
func (m *myint) print() {
	fmt.Println(*m)
}

func TestCustom(t *testing.T) {
	//自定义类型，是一种新的类型,需要强制转换,且原始的类型也可以是自定义类型
	int32Type := int32(1)
	p := myint(int32Type)
	p.print()

}

/*
例中，tempString是string的别名，其本质上与string是同一个类型。类型别名只会在代码中存在，编译完成后不会有如tempString一样的类型别名。所以变量s的类型是string。
字符类型中的byte和rune就是类型别名：

type byte = uint8
type rune = int32
类型别名这个功能非常有用，鉴于go中有些类型写起来非常繁琐，比如json相关的操作中，经常用到map[string]interface {}这种类型，写起来是不是很繁琐，没关系，给它起个简单的别名!这样用起来爽多了。

type strMap2Any = map[string]interface {}*/

type tempString = string

func TestAlisa(t *testing.T) {
	//自定义类型不用强装
	var stringType string

	stringType = "foo"

	var s tempString
	s = stringType

	fmt.Println(s)        // 我是s
	fmt.Printf("%T\n", s) // string
}

type alisaMyInt = myint

func (alisaMyInt) printValue() {

}

//这段代码编译会报错，自定义类型只能在原始类型所在的包中定义方法
//func (tempString) print(){
//
//}
func TestDefer(t *testing.T) {
	err := func1()
	if err != nil {
		fmt.Println("after func err =", err)
	}

}
func func1() (err error) {
	//defer 在return 后执行
	defer func() {
		if err != nil {
			log.Println("defer err =", err)
			err = nil
		}
	}()
	return errors.New("not found")
}

//interface object
type Mover interface {
	move()
}

type dog struct{}

//指针接收 对象是值类型不能实现接口。
func (d *dog) move() {
	fmt.Println("狗会动")
}
func TestImp(t *testing.T) {
	//var x Mover
	//var wangcai = dog{} // 旺财是dog类型
	//x = wangcai         // x不可以接收dog类型
	//var fugui = &dog{}  // 富贵是*dog类型
	//x = fugui           // x可以接收*dog类型
}

type Net interface {
	net.Conn
	Foo1()
}
type NetImp struct {
}

func (netImp *NetImp) Foo1() {

}

// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (netImp *NetImp) Read(b []byte) (n int, err error) {
	return 0, nil
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (netImp *NetImp) Write(b []byte) (n int, err error) {
	return 0, nil
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (netImp *NetImp) Close() error {
	return nil
}

// LocalAddr returns the local network address, if known.
func (netImp *NetImp) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr returns the remote network address, if known.
func (netImp *NetImp) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail instead of blocking. The deadline applies to all future
// and pending I/O, not just the immediately following call to
// Read or Write. After a deadline has been exceeded, the
// connection can be refreshed by setting a deadline in the future.
//
// If the deadline is exceeded a call to Read or Write or to other
// I/O methods will return an error that wraps os.ErrDeadlineExceeded.
// This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
// The error's Timeout method will return true, but note that there
// are other possible errors for which the Timeout method will
// return true even if the deadline has not been exceeded.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (netImp *NetImp) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (netImp *NetImp) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (netImp *NetImp) SetWriteDeadline(t time.Time) error {
	return nil
}

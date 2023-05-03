package standard

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"testing"
)

func TestBytes(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0, 4))
	buffer.Write([]byte{'a', 'c', 'd'})
	buf := make([]byte, 4)
	_, err := io.ReadFull(buffer, buf)
	if err != nil {
		log.Println(err)
	}
	log.Println(buffer.Bytes())
}
func TestStringByte(t *testing.T) {
	i := "112121212"
	log.Println([]byte(i))
}
func TestBinary(t *testing.T) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 258)
	//binary.BigEndian.Uint64()
	log.Println(buf)

}

//256 uint64 大端高位在前低地址
//八个字节
//0       0          0      0          0       0         1        0
//0000000 00000000 00000000 00000000 00000000 00000000 00000001 00000000

func TestBinary1(t *testing.T) {
	//序列化
	var dataA uint64 = 6010
	var buffer bytes.Buffer
	err1 := binary.Write(&buffer, binary.BigEndian, &dataA)
	if err1 != nil {
		log.Panic(err1)
	}
	byteA := buffer.Bytes()
	fmt.Println("序列化后：", byteA)

	//反序列化
	var dataB uint64
	var byteB []byte = byteA
	err2 := binary.Read(bytes.NewReader(byteB), binary.BigEndian, &dataB)
	if err2 != nil {
		log.Panic(err2)
	}
	fmt.Println("反序列化后：", dataB)
}
func TestTruncate1(t *testing.T) {
	//序列化
	buf := bytes.NewBuffer(make([]byte, 0, 1023))
	buf.Write([]byte{1, 2, 3})
	buf.Truncate(2)
	buf.Reset()
	log.Println(buf.Bytes())
}

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"testing"
)

const REQ_GENRANDOM_APP = 0x2008

func main1() {

	//建立TCP连接
	conn, err := net.Dial("tcp", "192.168.6.130:9166")
	if err != nil {
		log.Printf("%d: dial error: %s", 1, err)
		return
	}
	log.Println(1, ":connect to server ok")
	defer conn.Close()

	var sendData = make([]byte, 28)

	//以大端序写入总包长（）
	binary.BigEndian.PutUint32(sendData, uint32(24))
	//写入appid字符串
	copy(sendData[4:], []byte("1"))
	copy(sendData[8:], "####")
	//写入消息类型
	binary.BigEndian.PutUint32(sendData[12:], uint32(REQ_GENRANDOM_APP))
	binary.BigEndian.PutUint32(sendData[16:], uint32(0))
	binary.BigEndian.PutUint32(sendData[20:], uint32(4))
	//写入随机数长度
	binary.BigEndian.PutUint32(sendData[24:], uint32(16))
	fmt.Println("sendData", sendData)

	//发送数据
	rv, _ := conn.Write(sendData)
	fmt.Println("sendLen: ", rv)

	var recvData = make([]byte, 128)
	//接收数据
	recvLen, _ := conn.Read(recvData)
	fmt.Println("recvLen: ", recvLen)
	fmt.Println("recvData", recvData[:recvLen])

	return

}
func TestMashal(t *testing.T) {

}

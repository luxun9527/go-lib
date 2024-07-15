package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	// 连接到服务器
	conn, err := net.Dial("tcp", "192.168.2.159:8080")
	if err != nil {
		log.Fatalf("Error connecting to server: %s", err)
	}
	defer conn.Close()

	log.Println("Connected to server. Type messages and press enter to send.")

	// 创建一个新的读写器
	reader := bufio.NewReader(os.Stdin)

	for {
		log.Print("Enter message: ")
		// 从标准输入读取一行数据
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading from stdin: %s", err)
		}

		// 发送数据到服务器
		_, err = conn.Write([]byte(text))
		if err != nil {
			log.Fatalf("Error writing to connection: %s", err)
		}

		// 接收来自服务器的响应
		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			log.Fatalf("Error reading from connection: %s", err)
		}

		log.Printf("Received from server: %s", string(response[:n]))
	}
}

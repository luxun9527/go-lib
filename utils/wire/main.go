package main

import (
	"log"
	"net"
	"os"
	"time"
)

func server() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal("server", err)
		os.Exit(1)
	}
	data := make([]byte, 1)
	if _, err := conn.Read(data); err != nil {
		log.Fatal("server", err)
	}
	conn.Close()
}
func client() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("client", err)
	}
	if _, err := conn.Write([]byte("ab")); err != nil {
		log.Printf("client: %v", err)
	}
	time.Sleep(1 * time.Second) // wait for close on the server side

	if _, err := conn.Write([]byte("b")); err != nil {
		log.Printf("client: %v", err)
	}
	os.SyscallError{}
}

func main() {
	go server()
	time.Sleep(3 * time.Second) // wait for server to run

	client()
}

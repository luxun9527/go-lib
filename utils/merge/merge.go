package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func main(){
	const (
		msg = "message"
		service = "service"
		end="}"
	)
	fs, err := os.Open("E:\\GoCode\\go-lib\\utils\\merge\\mutl_message.proto")
	if err!=nil{
		log.Fatal(err)
	}
	target, err := os.OpenFile("E:\\GoCode\\go-lib\\utils\\merge\\mult.proto",os.O_CREATE|os.O_APPEND|os.O_RDWR,0755)
	if err!=nil{
		log.Fatal(err)
	}
	r := bufio.NewReader(fs)
	isStart :=false
	for {
		// ReadLine is a low-level line-reading primitive.
		// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		b, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		b = append(b, '\n')
		log.Println(string(b))
		if err != nil {

		}

		if  bytes.Contains(b, []byte(end)) || isStart{
			isStart = false
			target.Write(b)
			continue
		}
		if  bytes.Contains(b, []byte(msg)) || bytes.Contains(b, []byte(service)){
			isStart = true
			target.Write(b)
		}

	}

}



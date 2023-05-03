package main

import (
	"log"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("/bin/bash", "-c", "/home/deng/smb/go-lib/shell/shelldemo/stdin/reader/reader")
	pipe, err := cmd.StdinPipe()
	defer pipe.Close()
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(time.Second)
			pipe.Write([]byte("test\n"))
		}
	}()
	go func() {

		buf := make([]byte, 1024)
		for {
			n, err := stdoutPipe.Read(buf)
			if err != nil {
				log.Println("err", err)
			}
			log.Print(string(buf[:n]))
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal("err", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Println("err", err)
	}
}

package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	go func() {
		time.Sleep(time.Second)
		file, err := os.OpenFile("/home/deng/cloudsync.log", os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		for {
			time.Sleep(time.Second)
			file.WriteString(time.Now().String() + "asdaasdf\n")
		}
	}()

	cmd := exec.Command("/bin/bash", "-c", "tail -f /home/deng/cloudsync.log")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("pipe", err)
	}
	go func() {
		time.Sleep(time.Second * 3)
		log.Println("pid", cmd.Process.Pid)
		if err := cmd.Process.Signal(syscall.SIGSTOP); err != nil {
			//log.Println("signal err", err)
		}
	}()
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := pipe.Read(buf)
			if err != nil && err != io.EOF {
				log.Println("read err", err)
				return
			}
			if err == io.EOF {
				return
			}
			log.Println("read message", string(buf[:n]))
		}
	}()
	if err := cmd.Run(); err != nil {
		log.Println("err", err)
	}
	time.Sleep(time.Hour)
}

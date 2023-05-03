package main

import (
	"context"
	"io"
	"log"
	"os/exec"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sendCmd := exec.CommandContext(ctx, "/bin/bash", "-c", `parted /dev/sdg resizepart 1 2700MB`)
	stdinPipe, err := sendCmd.StdinPipe()
	stderr, err := sendCmd.StderrPipe()
	if err != nil {
		log.Println("StdinPipe", err)
	}
	stdoutPipe, err := sendCmd.StdoutPipe()
	go func() {
		defer stdoutPipe.Close()
		r := make([]byte, 4096)
		for {
			n, err := stdoutPipe.Read(r)
			if err == io.EOF {
				return
			}
			log.Println("StdoutPipe", string(r[:n]))
		}
	}()
	go func() {
		defer stderr.Close()
		r := make([]byte, 4096)
		for {
			n, err := stderr.Read(r)
			if err == io.EOF {
				return
			}
			log.Println("stderr", string(r[:n]))
		}
	}()
	defer stdinPipe.Close()

	if err := sendCmd.Start(); err != nil {
		log.Println("err", err)
		return
	}
	stdinPipe.Write([]byte("Y"))
	stdinPipe.Write([]byte("\n"))
	stdinPipe.Write([]byte("2800MB"))
	stdinPipe.Write([]byte("\n"))
	if err := sendCmd.Wait(); err != nil {
		log.Println("Wait", err)
		return
	}

}

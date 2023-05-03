package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", `/home/deng/smb/go-lib/shell/shelldemo/wait/output/output`)
	outPipe, err := cmd.StdoutPipe()
	defer outPipe.Close()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	br := bufio.NewReader(outPipe)
	go func() {
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			fmt.Println(string(a))
		}
	}()

	//fmt.Fscan()
	//var data string
	//for {
	//	_, err := fmt.Fscanln(outPipe, &data)
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Print(data)
	//}
	//buf := make([]byte, 2034)
	//go func() {
	//	for {
	//		n, err := outPipe.Read(buf)
	//		if err == io.EOF {
	//			break
	//		}
	//		if err != nil {
	//			panic(err)
	//
	//		}
	//		log.Print(string(buf[:n]))
	//	}
	//}()

	cmd.Wait()

}

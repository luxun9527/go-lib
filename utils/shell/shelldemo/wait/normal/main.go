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
	if err != nil {
		panic(err)
	}
	defer outPipe.Close()
	go func() {
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
	}()
	if err := cmd.Run(); err != nil {
		panic(err)
	}

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

	//cmd.Wait()

}

package main

import (
	"bytes"
	"log"
	"os/exec"
)

func main() {

	echoStderr()
}
func echoErrorCmd() {
	//1、错误命令
	cmd := exec.Command("/bin/bash", "-c", "ls /")
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stdout = &stderrBuf
	if err := cmd.Run(); err != nil {
		log.Println("err", err)
		return
	}
	log.Println("stdout", string(stdoutBuf.Bytes()))
	log.Println("stderr", string(stderrBuf.Bytes()))
}
func echoStderr() {
	cmd := exec.Command("/bin/bash", "-c", "/home/deng/smb/go-lib/shell/shelldemo/wait/echostderr/echostderr")
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if err := cmd.Run(); err != nil {
		log.Println("err", err)
		return
	}
	log.Println("stdout", string(stdoutBuf.Bytes()))
	log.Println("stderr", string(stderrBuf.Bytes()))

}

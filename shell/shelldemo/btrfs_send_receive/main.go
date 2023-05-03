package main

import (
	"context"
	"io"
	"os/exec"
)

func main() {
	// 必须要使用root
	//1 读取文件，转换为流
	//2、拷贝文件到输出流
	//3、存储
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sendCmd := exec.CommandContext(ctx, "/bin/bash", "-c", `btrfs send /btrfsdir/src`)
	sendPipe, err := sendCmd.StdoutPipe()
	defer sendPipe.Close()
	if err != nil {
		panic(err)
	}
	if err := sendCmd.Start(); err != nil {
		panic(err)
	}
	//fmt.Fscanln(sendPipe)
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", `btrfs receive  /btrfsdir/dist`)

	receivePipe, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	if _, err := io.Copy(receivePipe, sendPipe); err != nil {
		panic(err)
	}

	//receivePipe.Close()
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	//
	//cmd := exec.Command("echo", "-n", `{"Name": "Bob", "Age": 32}`)
	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := cmd.Start(); err != nil {
	//	log.Fatal(err)
	//}
	//var person struct {
	//	Name string
	//	Age  int
	//}
	//if err := json.NewDecoder(stdout).Decode(&person); err != nil {
	//	log.Fatal(err)
	//}
	//if err := cmd.Wait(); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%s is %d years old\n", person.Name, person.Age)
}

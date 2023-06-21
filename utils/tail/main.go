package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hpcloud/tail"
)

func main() {
	fileName := "S:\\go-lib\\tail\\my.log"
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	go func() {
		time.Sleep(time.Second * 3)
		for {
			file, err := os.OpenFile("S:\\go-lib\\tail\\my.log", os.O_RDWR|os.O_APPEND, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second)
			file.WriteString(" *           CloudSync_TOS_APP_2.4.14_x86_64.tpk:100% /48.583Mi, 1.074Mi/s, 0sTransferred:         51.343 MiB / 51.343 MiB, 100%, 1.008 MiB/s, ETA 0s\nTransferred:            8 / 9, 89%\nElapsed time:      1m42.2s\nTransferring:")
		}

	}()
	for {
		line, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("line:", line.Text)
	}
}

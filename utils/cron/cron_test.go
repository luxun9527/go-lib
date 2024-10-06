package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"testing"
	"time"
)

// https://help.aliyun.com/document_detail/133509.html
func TestCron(t *testing.T) {
	c := cron.New(
		cron.WithLogger(
			cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
		cron.WithLocation(time.Local),
		cron.WithSeconds(), //开启秒一级
	)
	//if _, err := c.AddFunc("@every 1s", func() {
	//	fmt.Println("hello world")
	//}); err != nil {
	//	log.Panicf("cron add func failed, err:%v\n", err)
	//}
	//每两秒执行一次
	if _, err := c.AddFunc("0/2 * * * * ?", func() {
		fmt.Println("hello world")
	}); err != nil {
		log.Panicf("cron add func failed, err:%v\n", err)
	}
	c.Start()
	time.Sleep(5 * time.Second)
}

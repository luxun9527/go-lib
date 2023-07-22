package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"testing"
	"time"
)

func TestCron(t *testing.T) {
	c := cron.New(
		cron.WithLogger(
			cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
		cron.WithSeconds(),
	)
	c.AddFunc("@every 1s", func() {
		fmt.Println("hello world")
	})
	c.Start()
	time.Sleep(5 * time.Second)
}

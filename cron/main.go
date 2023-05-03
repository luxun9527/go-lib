package main

import (
	"github.com/robfig/cron/v3"
	"gitlab.local/golibrary/crontable"
	"log"
	"time"
)

func main() {
}

type Task struct {
}

func (t *Task) Run() {
	log.Println("111111")
}
func cronTable() {
	c := crontable.New()
	c.AddJob("0/30 * * * * *", func() {
		log.Println("test")
	})
	c.Start()

	time.Sleep(time.Hour)
}
func example() {
	tString := time.Now().Format("GMT-07-2006.01.02-15.04.05")
	log.Println(tString)
	c := cron.New(
		cron.WithSeconds(),
	)

	//c.AddFunc("* * * * * *", func() {
	//	fmt.Println("tick every 1 second")
	//})
	time.Now().Minute()
	c.Start()
	var task Task
	c.AddJob("0 */1 * * * *", &task)
}
func Schedule() {
	c := cron.New(
		cron.WithSeconds(),
	)

	c.Start()
}

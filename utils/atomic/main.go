package main

import (
	"go.uber.org/atomic"
	"log"
)

func main() {
	var i atomic.Int64
	i.Inc()
	i.Add(1)
	i.Add(1)
	log.Println(i.Load())
}

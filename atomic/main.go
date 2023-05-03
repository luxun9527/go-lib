package main

import (
	"go.uber.org/atomic"
	"log"
)

func main() {
	var i atomic.Int64
	i.Add(1)
	i.Add(1)
	log.Println(i.Load())
}

package main

import (
	cb "github.com/emirpasic/gods/queues/circularbuffer"
	"github.com/emirpasic/gods/stacks/arraystack"
	"testing"
)

func TestQueue(t *testing.T) {

	queue := cb.New(5)
	queue.Enqueue(1)
	queue.Enqueue(2)
	value, ok := queue.Dequeue()

	if ok {
		t.Log(value)
	}

	an := arraystack.New()
	
	for _, v := range an.Values() {
		t.Log(v)
	}
}

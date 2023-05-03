package main

import (
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"

	"testing"
)

//go test -v benchmark_test.go -benchtime=5s  -bench=BenchmarkList$ -count=3
func BenchmarkList(b *testing.B) {
	list := fillList()
	iterator := list.Iterator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for iterator.Next() {
			_, _ = iterator.Index(), iterator.Value()
		}
	}
}

const N = 1000

func fillList() *sll.List {
	list := sll.New()
	for i := 0; i < N; i++ {
		list.Add("asdf")
	}
	return list
}

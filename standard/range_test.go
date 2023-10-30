package standard

import (
	"log"
	"testing"
)

func TestRange(t *testing.T) {
	type student struct {
		name string
	}
	//range是值覆盖的方式。
	students,target :=make([]student,0,10) ,make(map[string]*student, 10)
	students = append(students, student{name: "1"},student{name: "2"})
	for _,v := range students {
		target[v.name]=&v
		log.Printf("%p",&v)
	}
	for _,v := range target {
		log.Printf("%+v",v)
	}
}

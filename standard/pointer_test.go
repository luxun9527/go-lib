package standard

import (
	"log"
	"testing"
)

func TestPointer(t *testing.T) {
	type Person struct {
		Name string
		Age  string
	}
	person := &Person{
		Name: "zhangsan",
		Age:  "12",
	}
	p := *person
	p.Name = "lisi"
	log.Printf("%+v", p)
	log.Printf("%+v", person)
	var i = 2
	log.Println(-i)
}

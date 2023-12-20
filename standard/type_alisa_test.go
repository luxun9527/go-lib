package standard

import (
	"log"
)

type MyInt int32

type IntAlisa = int32

func (MyInt) New() {

}

//func (IntAlsa) New() { 错误
//
//}

type P interface {
	Print1()
	Print2()
}


type P1 struct {
	Name string
}
func(p P1) Print1(){
	log.Println("p1",p.Name)
}
func(p P1) Print2(){
	log.Println("p1",p.Name)
}
type P2 P1

func (p P2)Print1(){
	log.Println("p2",p.Name)
}
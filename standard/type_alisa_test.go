package standard

import (
	"log"
	"testing"
)

type MyInt int32

type IntAlisa = int32

func (MyInt) New() {

}

//func (IntAlsa) New() { 错误
//
//}
func TestAlisa(t *testing.T) {
	findThirdValue()
}
func findThirdValue() {
	m, index, third := []int32{1, 3, 5, 90, 2, 9}, 0, int32(0)
	log.Printf("%p", m)
	for i := 0; i < 3; i++ {
		max := int32(0)
		for k, v := range m {
			if v > max {
				max = v
				index = k
			}
		}
		if index != len(m)-1 {

			m = append(m[:index], m[index+1:]...)
			log.Printf("%p", m)
		} else {
			m = append(m[:index])
		}
		third = max
	}
	log.Println(third)
}

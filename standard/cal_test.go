package standard

import (
	"log"
	"testing"
)

func TestCal(t *testing.T) {

	log.Println(3000/1e2)
	//0 0 0 0 0 0 1 0
	//0 0 0 1 0 0 0 0
	// 2乘以2的3次方 16
	log.Println(2 <<3)
	//0 0 0 0 0 1 0 0
	//0 0 0 0 0 0 1 0
	// 4除以2的1次方 2
	log.Println(4 >>1)


}

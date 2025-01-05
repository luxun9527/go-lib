package standard

import (
	"errors"
	"log"
	"testing"
)

func f1() (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	err = errors.New("test")
	if err != nil {
		return err
	}
	return errors.New("test error")
}
func f2() {}
func TestDefer(t *testing.T) {
	f1()
}

func TestDefer2(t *testing.T) {
	var i int = 0
	defer func(j int) {
		log.Printf("func1 %v", j)
	}(i)
	i++
	defer func(k int) {
		log.Printf("func2 %v", k)
	}(i)

	panic("err")
}

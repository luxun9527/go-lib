package standard

import (
	"errors"
	"log"
	"testing"
)

func f1()(err error){
	defer func() {
		if err!=nil{
			log.Println(err)
		}
	}()

	err = errors.New("test")
	if err!=nil{
		return err
	}
	return errors.New("test error")
}
func f2(){}
func TestRun(t *testing.T) {
	f1()
}
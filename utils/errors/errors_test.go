package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
	"log"
	"testing"
)

var (
	NotFound = MyError{
		Code:    404,
		Message: "Not Found",
	}
	InternalError = MyError{
		Code:    500,
		Message: "Internal Error",
	}
)

type MyError struct {
	Code    int32
	Message string
}

func (e MyError) Error() string {
	return e.Message
}
func FromError(err error) MyError {
	myError, ok := err.(MyError)
	if ok {
		return myError
	}
	return InternalError
}

func TestErrors(t *testing.T) {
	err := GetError()
	fromError := FromError(err)
	log.Println(fromError)
	err = errors.Wrap(errors.New("not found"), "warp message")
	log.Println(err)

}
func Wrap(err error,message string)error{
	s, _ := status.FromError(err)
	msg := fmt.Sprintf("%v:%v",s.Message(),message)
	return status.New(s.Code(),msg).Err()
}

func GetError() error {
	return NotFound
}

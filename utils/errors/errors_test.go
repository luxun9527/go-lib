package errors

import (
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
}
func GetError() error {
	return NotFound
}

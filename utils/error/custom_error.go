package error

import (
	"errors"
	"fmt"
)

type CustomError struct {
	Code    int32
	Message string
}

// 自定义错误类型需要实现error接口
func (e CustomError) Error() string {
	return e.Message
}
func (e CustomError) String() string {
	return fmt.Sprintf(`{code:"%v",message:"%v"}`, e.Code, e.Message)
}

// errors.Is()会调用这个方法来判断是否是同一个错误类型
func (e CustomError) Is(err error) bool {
	var customError CustomError
	if ok := errors.As(err, &customError); !ok {
		return false
	}
	return e.Code == customError.Code

}

// 提取指定类型的错误，判断包装的 error 链中，某一个 error 的类型是否与 target 相同，并提取第一个符合目标类型的错误的值，将其赋值给 target。
func (e CustomError) As(target interface{}) bool {

	t, ok := target.(**E)
	if ok {
		*t = &E{}
	}
	return true
}

// 解析错误，获取原始错误
func (e CustomError) Unwrap() error {
	return e
}

func NewCustomError(code int32, message string) CustomError {
	return CustomError{
		Code:    code,
		Message: message,
	}

}

type E struct {
}

func (E) Error() string {
	return "E type error"
}

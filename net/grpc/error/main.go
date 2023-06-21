package main

import (
	"errors"
	"github.com/gin-gonic/gin"
)

/*
	错误返回设计
	1、错误码能够体现出相对详细的错误，message暴露出更少的信息，详细的错误信息一般在日志中。
	2、定义在同一的地方，所有的服务去依赖。
	3、兼顾易用性，提供易用灵活的api。
*/

type Code int32

// 通用错误
const (
	// Unknown 未知错误
	Unknown Code = 100000
	// ParamValidateFailed 参数校验失败
	ParamValidateFailed Code = 110000
	// HeadValidateFailed 通用请求头校验失败
	HeadValidateFailed Code = 111000
	// ExecSqlFailed Sql执行失败
	ExecSqlFailed Code = 120000
	// RecordNotFound 在指定条件下查找有记录没找到。
	RecordNotFound Code = 121000
	// RedisFailed 使用redis错误
	RedisFailed Code = 130000
	// MongoFailed 使用mongo错误 。
	MongoFailed Code = 131000
	// KafkaFailed kafka错误 。
	KafkaFailed Code = 132000
	// EtcdFailed 使用Etcd错误 。
	EtcdFailed Code = 133000
	// AuthFailed 认证失败
	AuthFailed Code = 140000
	// Timeout 超时
	Timeout Code = 150000
)

type Error struct {
	Code    Code   `json:"code"`    // 业务码
	Message string `json:"message"` // 描述信息
}

func (e *Error) Error() string {
	return e.Message
}
func Convert(err error) Error {
	e, ok := err.(*Error)
	if ok {
		return *e
	}
	return NewError(Unknown, err.Error())
}

func (e *Error) Is(err error) bool {
	return errors.Is(e, err)
}
func (e *Error) WithMessage(message string) {
	e.Message = message
}
func NewError(code Code, message string) Error {
	return Error{code, message}
}
func NewCodeError(code Code) Error {
	return Error{Code: code}
}

func main() {
	r := gin.New()

	r.GET("/test", func(c *gin.Context) {

		//c.JSON(200)
	})
	r.Run(":9090")
}

//go:generate errgen -p u.code.go

package main

// CommonCodeInit 定义公共错误码的起始位置

const (
	InternalCode = iota + 1
	RedisErrCode //redis错误
	ParamValidateErrCode //参数校验失败
	RecordNotFoundErrCode//记录未找到
	DuplicateDataErrCode//重复数据
)




package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



const (
	UnAuthorizedCode codes.Code = iota+403 //403未授权
	NotFoundCode  //404未找到
)


var (
	NotFound = status.Error(UnAuthorizedCode,"未找到")
)



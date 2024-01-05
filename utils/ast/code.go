package main

import (
	"google.golang.org/grpc/codes"
)

const (
	UnAuthorizedCode codes.Code = iota + 403 //403未授权
	NotFoundCode                             //404未找到
)

package error

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Code codes.Code

func (c Code) Translate(lang string) string {
	d, ok := _m[lang]
	if !ok {
		d = _m["zh-CN"] //默认中文
	}
	result, ok := d[c]
	if !ok {
		return "服务器繁忙,请稍后再试"
	}
	return result
}

func (c Code) Error(msg string) error {
	return status.New(codes.Code(c), msg).Err()
}

var _m = map[string]map[Code]string{
	"zh-CN": {
		UserNotFoundCode: "用户不存在",
		InternalCode:     "服务器繁忙,请稍后再试",
	},
}

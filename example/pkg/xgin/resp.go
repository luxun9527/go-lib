package xgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go-lib/example/pkg/i18n"
	"google.golang.org/grpc/status"
	"net/http"
)

type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

var (
	Empty = struct{}{}
)

func FailWithLangError(c *gin.Context, err error) {
	ResponseWithLang(c, struct{}{}, err)
}
func FailWithLang(c *gin.Context) {
	FailWithLangError(c, fmt.Errorf("unknown error"))
}
func ResponseWithLang(c *gin.Context, resp interface{}, err error) {
	lang := c.GetHeader("lang")
	if err != nil {
		code := status.Code(err)
		msg := i18n.Translate(lang, cast.ToString(uint32(code)))
		Result(cast.ToInt(uint32(code)), Empty, msg, c)
	} else {
		Result(SUCCESS, resp, "success", c)
	}
}

func Response(c *gin.Context, resp interface{}, err error) {
	if err != nil {
		r, _ := status.FromError(err)
		Result(cast.ToInt(r.Code()), Empty, r.Message(), c)
	} else {
		Result(SUCCESS, resp, "success", c)
	}
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, CommonResponse{
		code,
		msg,
		data,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "查询成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

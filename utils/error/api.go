package error

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

func InitApiSrv() {
	engine := gin.Default()

	var d DemoApi
	engine.GET("/error", d.Demo)
	log.Printf("start http server :9999")
	if err := engine.Run(":9999"); err != nil {
		log.Panicf("%v", err)
	}

}

type DemoApi struct {
}

func (h *DemoApi) Demo(c *gin.Context) {
	var d DemoService
	err := d.Demo()
	Response(c, Nil, err)
}

var Nil = struct{}{}

// Response 在api层统一处理错误
func Response(c *gin.Context, data interface{}, err error) {
	if err != nil {
		header := c.GetHeader("language")
		if header == "" {
			header = "zh-CN"
		}
		e, ok := status.FromError(err)
		code := Code(e.Code())
		msg := code.Translate(header)
		if !ok || uint32(code) < uint32(CommonCodeInit) {
			msg = InternalCode.Translate(header)
		}
		c.JSON(200, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})

		return
	}
}

// DemoService api层的业务逻辑
type DemoService struct {
}

func (h *DemoService) Demo() error {
	conn, err := grpc.Dial("127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	cli := NewErrorsClient(conn)
	if _, err := cli.Demo(context.Background(), &Empty{}); err != nil {
		//不用关心，调用的grpc接口返回的是什么直接返回，api层最后统一处理
		log.Printf("call grpc demo failed %v", err)
		return err
	}
	return nil
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(Cors())
	//绑定get参数
	route.GET("/api/test/time", func(c *gin.Context) {
		h := gin.H{"data": time.Now().Format(time.DateTime), "code": 200}
		c.JSON(200, h)
	})
	fmt.Scan(1)
	route.Run(":8080")
}

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	//绑定get参数
	route.GET("/test", func(c *gin.Context) {

		select {
		}
	})
	route.Run(":9090")
}


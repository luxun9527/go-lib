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
	route.GET("/api/omcenter/", func(c *gin.Context) {
		//c.ShouldBindUri()
		hs := make([]gin.H, 0, 10)
	//	l := []gin.H{{"name":"test1"},{"name":"test2"}}
		//data, _ := json.Marshal(l)
		c.JSON(200,hs)
	})
	route.Run(":8200")
}


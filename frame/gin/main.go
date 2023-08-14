package main

import (
	"github.com/gin-gonic/gin"
	"log"
)




func main() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Recovery())
	//绑定get参数
	route.GET("/test", func(c *gin.Context) {
		log.Panic("test")
	})
	route.Run(":9090")
}


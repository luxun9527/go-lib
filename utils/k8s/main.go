package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

func main() {

	route := gin.Default()
	route.GET("/api/time", func(c *gin.Context) {
		h := gin.H{"time": time.Now().Format(time.DateTime), "code": 200}
		c.JSON(200, h)
	})
	route.Run(":8089")
}

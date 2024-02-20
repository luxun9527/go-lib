package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.New()
	engine.POST("/auth/login", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "data": "testet"})
	})
}

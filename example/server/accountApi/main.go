package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Run(":8080")
}

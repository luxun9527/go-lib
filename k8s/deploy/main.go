package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {

	e := gin.Default()
	e.GET("/getTime", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 200, "data": time.Now().Unix()})
	})
	if err := e.Run(":10001"); err != nil {
		log.Panicf("startup service failed, err:%v\n", err)
	}
}

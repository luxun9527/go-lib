package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(Cors())
	//绑定get参数
	route.POST("/upload", func(c *gin.Context) {
		_, headers, err := c.Request.FormFile("file")
		if err != nil {
			log.Printf("Error when try to get file: %v", err)
		}

		if err := c.SaveUploadedFile(headers, "./video/"+cast.ToString(time.Now().UnixNano())+"/"+headers.Filename); err != nil {
			return
		}
		c.String(http.StatusOK, headers.Filename)
	})
	route.Run(":8080")
}

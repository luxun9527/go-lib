package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
func main() {
	// 初始化 Gin
	r := gin.Default()

	// 配置 CORS 中间件
	r.Use(Cors())

	// 定义一个路由
	r.POST("/api/data", func(c *gin.Context) {
		// 设置共享的 Cookie
		host, _, err := net.SplitHostPort(c.Request.Host)
		if err != nil {
			host = c.Request.Host
		}
		host = strings.TrimLeft(host, "api.")
		if net.ParseIP(host) != nil {
			c.SetCookie("x-token", "token", 3600, "/", "", false, false)
		} else {
			c.SetCookie("x-token", "test", 3600, "/", host, false, false)
		}
		// 返回响应
		c.JSON(200, gin.H{
			"message": "Hello from API!",
		})
	})

	// 定义一个路由
	r.POST("/api/getCookie", func(c *gin.Context) {
		// 设置共享的 Cookie
		cook, err := c.Cookie("x-token")
		if err != nil {
			log.Printf("%v", err)
			return
		}
		log.Println(cook)
		// 返回响应
		c.JSON(200, gin.H{
			"message": "Hello from API!",
		})
	})
	// 启动服务器
	r.Run(":9898")
}

package cors

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
)

func TestCors(t *testing.T) {
	engine := gin.New()
	routes := engine.Use(Cors())
	routes.POST("/login", func(c *gin.Context) {
		log.Println("aaaa")
	})
	engine.Run(":8888")
}

//https://www.ruanyifeng.com/blog/2016/04/cors.html
//https://juejin.cn/post/6844903850684465159
//https://segmentfault.com/q/1010000040723727
//https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CORS
//https://cloud.tencent.com/developer/article/2131856
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

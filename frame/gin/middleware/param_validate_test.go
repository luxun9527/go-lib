package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"testing"
)

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// https://www.kancloud.cn/shuangdeyu/gin_book/949432
func TestParamValidate(t *testing.T) {
	router := gin.Default()
	router.POST("/login", func(c *gin.Context) {
		// you can bind multipart form with explicit binding declaration:
		// c.ShouldBindWith(&form, binding.Form)
		// or you can simply use autobinding with ShouldBind method:
		var form LoginForm
		// in this case proper binding will be automatically selected
		if err := c.ShouldBind(&form); err != nil {
			log.Printf("error during binding: %v", err)
			c.JSON(400, gin.H{"status": "invalid request"})
		}
		if form.User == "user" && form.Password == "password" {
			c.JSON(200, gin.H{"status": "you are logged in"})
		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
		}

	})
	router.Run(":8080")
}

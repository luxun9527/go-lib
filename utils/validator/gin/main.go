package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Username string `form:"username" binding:"required"`
	Password string ` form:"password" binding:"required"`
}

func main() {
	r := gin.New()
	translator, _ := NewTranslator(binding.Validator.Engine().(*validator.Validate))
	//binding.Validator.Engine()
	r.POST("/addUser", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {

			msg := translator.TranslateFirst("zh", err)
			c.JSON(200, gin.H{"success": 200, "message": msg})
			return
		}
		c.JSON(200, gin.H{"success": 200, "message": ""})

	})
	r.Run(":9999")
	/**

	 */
}

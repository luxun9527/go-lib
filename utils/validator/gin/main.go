package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Auth struct {
	AuthName string `gorm:"column:auth_name;type:varchar(50)" json:"authName" binding:"required"`
	Path     string `gorm:"column:path;type:varchar(100)" json:"path" binding:"required"`
}

func main() {
	r := gin.New()
	binding.Validator = new(DefaultValidator)
	binding.Validator.Engine()
	r.POST("/addAuth", func(c *gin.Context) {
		var auth Auth
		if err := c.ShouldBindJSON(&auth); err != nil {
			c.JSON(200, gin.H{"success": 200, "message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": 200, "message": ""})

	})
	r.Run(":9090")

}

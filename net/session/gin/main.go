package main

import (
	gsession "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	s := cookie.NewStore([]byte("secret"))

	r.Use(gsession.Sessions("mysession", s))
	r.GET("/setSession", func(c *gin.Context) {
		session := gsession.Default(c)
		session.Set("name", "zhangsan")
		session.Set("age", 24)
		session.Save()
		c.JSON(200, gin.H{"success": 200})
	})
	r.GET("/getSession", func(c *gin.Context) {
		session := gsession.Default(c)

		name := session.Get("name")
		log.Println("name = ", name)
		c.JSON(200, gin.H{"name": name})

	})
	r.GET("/deleteSession", func(c *gin.Context) {
		session := gsession.Default(c)
		res := session.Get("name")
		id := session.ID()
		log.Println(id)
		log.Println("res = ", res)
		session.Delete("name")
		session.Save()
	})

	r.Run(":8000")
}

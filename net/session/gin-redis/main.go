package main

import (
	gsession "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Println(err)
	}

	r.Use(gsession.Sessions("mysession", store))

	r.GET("/setRedisSession", func(c *gin.Context) {
		session := gsession.Default(c)
		session.Set("name", "zhangsan")
		if err := session.Save(); err != nil {
			log.Println("save failed", err)
		}
		c.JSON(200, gin.H{"success": "200"})
	})
	r.GET("/getRedisSession", func(c *gin.Context) {
		session := gsession.Default(c)
		name := session.Get("name")
		c.JSON(200, gin.H{"success": "200", "name": name})
	})
	r.GET("/deleteRedisSession", func(c *gin.Context) {
		session := gsession.Default(c)
		session.Delete("name")
		session.Save()
	})
	r.Run(":8000")
}

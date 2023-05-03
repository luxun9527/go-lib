package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/golang-jwt/jwt/v4/test"
	"jwt/model"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func main() {
	engine := gin.New()
	r := engine.Use(func(c *gin.Context) {
		var customerClaims model.CustomClaims
		token, err := request.ParseFromRequest(
			c.Request,
			request.AuthorizationHeaderExtractor,
			func(*jwt.Token) (interface{}, error) {
				return test.LoadRSAPublicKeyFromDisk("/Users/demg/personProject/go-lib/jwt/public.key"), nil
			},
			request.WithClaims(&customerClaims),
		)
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					c.JSON(200, gin.H{"message": err.Error()})
					c.Abort()
					return
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					c.JSON(200, gin.H{"message": err.Error()})
					c.Abort()
					return
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					c.JSON(200, gin.H{"message": err.Error()})
					c.Abort()
					return
				} else {
					c.JSON(200, gin.H{"message": err.Error()})
					c.Abort()
					return
				}
			} else {
				c.JSON(200, gin.H{"message": err.Error()})
				c.Abort()
				return
			}
		}
		if token != nil {
			if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
				c.Set("username", claims.UserName)
				c.Set("id", claims.UserId)
				c.Next()
			}
		}
	})
	r.GET("/get", func(c *gin.Context) {
		username, _ := c.Get("username")
		id, _ := c.Get("id")

		c.JSON(200, gin.H{"username": username, "id": id})
	})
	engine.Run(":10090")
}

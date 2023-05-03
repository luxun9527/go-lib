package model

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	UserName string
	UserId   int32
	jwt.RegisteredClaims
}

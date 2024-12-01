package xjwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/test"
	"go-lib/example/common/errs"
	"time"
)

var DefaultJwtConf *JwtConf

type JwtConf struct {
	SignKey           string
	RsaPrivateKeyPath string
	ValidTime         int64
}

var (
	InValidTokenErr = errs.TokenValidateFailed
	TokenExpiredErr = errs.TokenExpire
)

type CustomClaims[T any] struct {
	Extra   T
	jwtConf *JwtConf
	jwt.RegisteredClaims
}

func NewCustomClaims[T any](extra T, jwtConf ...*JwtConf) (*CustomClaims[T], error) {
	var jc *JwtConf
	if len(jwtConf) == 0 {
		jc = DefaultJwtConf
	} else {
		jc = jwtConf[0]
	}
	if jc == nil {
		return nil, fmt.Errorf("jwtConf is nil")
	}
	if jc.RsaPrivateKeyPath == "" && jc.SignKey == "" {
		return nil, fmt.Errorf("jwtConf.SignKey and jwtConf.RsaPrivateKeyPath are both empty")
	}
	return &CustomClaims[T]{
		Extra:   extra,
		jwtConf: jc,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(jc.ValidTime))),
		},
	}, nil
}

func (c *CustomClaims[T]) GenerateToken() (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	if c.jwtConf.SignKey != "" {
		return tokenClaims.SignedString([]byte(c.jwtConf.SignKey))
	} else {
		priKey := test.LoadRSAPrivateKeyFromDisk(c.jwtConf.RsaPrivateKeyPath)
		return tokenClaims.SignedString(priKey)
	}
}
func ParseToken[T any](tokenKey string, jwtConf ...*JwtConf) (*CustomClaims[T], error) {
	var jc *JwtConf
	if len(jwtConf) == 0 {
		jc = DefaultJwtConf
	} else {
		jc = jwtConf[0]
	}
	if jc == nil {
		return nil, fmt.Errorf("jwtConf is nil")
	}
	if jc.RsaPrivateKeyPath == "" && jc.SignKey == "" {
		return nil, fmt.Errorf("jwtConf.SignKey and jwtConf.RsaPrivateKeyPath are both empty")
	}

	var customClaims CustomClaims[T]
	token, err := jwt.ParseWithClaims(tokenKey, &customClaims, func(c *jwt.Token) (interface{}, error) {
		if jc.SignKey != "" {
			return []byte(jc.SignKey), nil
		} else {
			return test.LoadRSAPublicKeyFromDisk(jc.RsaPrivateKeyPath), nil
		}
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenExpiredErr
		}
		return nil, InValidTokenErr
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims[T]); ok && token.Valid {
			return claims, nil
		}
		return nil, InValidTokenErr

	} else {
		return nil, InValidTokenErr
	}
}

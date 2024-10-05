package test

//根据秘钥生成和解析jwt token
import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"testing"
	"time"
)

// https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
// https://pkg.go.dev/github.com/golang-jwt/jwt/v4#RegisteredClaims
type CustomClaims struct {
	UserName string
	UserId   int32
	jwt.RegisteredClaims
}

const (
	SignKey = "569ef72642be0fadd711d6a468d68ee1"
)

// ssh-keygen -t rsa -b 4096 -C "your_email@example.com
// 可以通过这个命令来生成秘钥
// 通过rsa 生成token
func GenerateToken() (string, error) {
	userInfo := CustomClaims{
		UserName: "zhangSan",
		UserId:   1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	//tokenClaims := jwt.NewWithClaims(jwt.SigningMethodRS256, userInfo)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, userInfo)
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	//priKey := test.LoadRSAPrivateKeyFromDisk("./private.key")
	token, err := tokenClaims.SignedString([]byte(SignKey))
	//token, err := tokenClaims.SignedString(priKey)
	return token, err
}

func TestParseToken(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		log.Panicf("generate token err: %v", err)
	}
	log.Printf("token: %s", token)
	claims, err := ParseToken(token)
	if err != nil {
		log.Printf("parse token err: %v", err)
	}
	log.Printf("claims: %v", claims)
}

var (
	InValidTokenErr = errors.New("无效的token")
	TokenExpiredErr = errors.New("token不在有效期")
)

func ParseToken(tokenKey string) (*CustomClaims, error) {
	var customClaims CustomClaims
	token, err := jwt.ParseWithClaims(tokenKey, &customClaims, func(c *jwt.Token) (interface{}, error) {
		//return test.LoadRSAPublicKeyFromDisk("./public.key"), nil
		return []byte(SignKey), nil
	})

	if err != nil {
		log.Printf("parse token err: %v", err)
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, TokenExpiredErr
		}
		return nil, InValidTokenErr
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, InValidTokenErr

	} else {
		return nil, InValidTokenErr
	}
}

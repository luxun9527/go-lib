package test

//根据秘钥生成和解析jwt token
import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/test"
	"jwt/model"
	"testing"

	"time"
)

//https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
//https://pkg.go.dev/github.com/golang-jwt/jwt/v4#RegisteredClaims

func GenerateToken() (string, error) {
	userInfo := model.CustomClaims{
		UserName: "zhangSan",
		UserId:   1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodRS256, userInfo)
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	priKey := test.LoadRSAPrivateKeyFromDisk("/Users/demg/personProject/go-lib/jwt/private.key")
	token, err := tokenClaims.SignedString(priKey)
	return token, err
}
func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		fmt.Println("generate token failed", err)
		return
	}
	fmt.Println("token = ", token)
}
func TestParseToken(t *testing.T) {
	token, err := GenerateToken()
	if err != nil {
		fmt.Println("generate token failed", err)
		return
	}
	fmt.Println("token = ", token)
	time.Sleep(time.Second * 2)
	if err := ParseToken(token); err != nil {
		fmt.Println("parseToken err", err)
	}
}
func ParseToken(token string) error {
	var customClaims model.CustomClaims
	result, err := jwt.ParseWithClaims(token, &customClaims, func(c *jwt.Token) (interface{}, error) {
		return test.LoadRSAPublicKeyFromDisk("/Users/demg/personProject/go-lib/jwt/public.key"), nil
	})

	if result.Valid {
		fmt.Println("You look nice today")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		fmt.Println("That's not even a token")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		fmt.Println("Timing is everything")
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}
	return nil
}

//func (*JWT) ParseToken(c *gin.Context) (*jwt2.CustomClaims, error) {
//	var customerClaims jwt2.CustomClaims
//	token, err := request.ParseFromRequest(
//		c.Request,
//		request.AuthorizationHeaderExtractor,
//		func(*jwt.Token) (interface{}, error) {
//			return test.LoadRSAPublicKeyFromDisk(global.GCONFIG.Jwt.PublicKeyPath), nil
//		},
//		request.WithClaims(&customerClaims),
//	)
//	if err != nil {
//		if ve, ok := err.(*jwt.ValidationError); ok {
//			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
//				return nil, TokenMalformed
//			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
//				// Token is expired
//				return nil, TokenExpired
//			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
//				return nil, TokenNotValidYet
//			} else {
//				return nil, TokenInvalid
//			}
//		}
//	}
//	if token != nil {
//		if claims, ok := token.Claims.(*jwt2.CustomClaims); ok && token.Valid {
//			return claims, nil
//		}
//		return nil, TokenInvalid
//
//	} else {
//		return nil, TokenInvalid
//
//	}
//
//}

package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-lib/example/common/errs"
	accountPb "go-lib/example/pb/account"
	"go-lib/example/pkg/xgin"
	"go-lib/example/server/accountApi/rpcClient"
)

func TokenValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			xgin.FailWithLangError(c, errs.TokenValidateFailed)
			c.Abort()
			return
		}
		resp, err := rpcClient.AccountClient.ValidateToken(context.Background(), &accountPb.ValidateTokenReq{Token: token})
		if err != nil {
			xgin.ResponseWithLang(c, resp, err)
			c.Abort()
			return
		}
		c.Set("accountId", resp.AccountId)
		c.Set("accountName", resp.AccountName)
		c.Next()
	}
}

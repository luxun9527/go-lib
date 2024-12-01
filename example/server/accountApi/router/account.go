package router

import (
	"github.com/gin-gonic/gin"
	"go-lib/example/server/accountApi/api/account"
	"go-lib/example/server/accountApi/middleware"
)

func initAccountRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("account")
	userRouter.Use(middleware.TokenValidator())
	{
		userRouter.GET("getAccountInfo", account.AccountApi.GetAccountInfo) // 获取用户
	}
}

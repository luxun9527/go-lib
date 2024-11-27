package router

import "github.com/gin-gonic/gin"

func initUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		userRouter.POST("createUser")        // 创建User
		userRouter.DELETE("deleteUser")      // 删除User
		userRouter.DELETE("deleteUserByIds") // 批量删除User
	}
}

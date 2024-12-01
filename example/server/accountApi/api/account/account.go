package account

import (
	"github.com/gin-gonic/gin"
	"go-lib/example/pkg/xgin"
	accountApiModel "go-lib/example/server/accountApi/model/account"
	accountService "go-lib/example/server/accountApi/service/account"
)

var AccountApi = &accountApi{}

type accountApi struct{}

func (*accountApi) GetAccountInfo(c *gin.Context) {
	accountId := c.GetString("accountId")
	resp, err := accountService.UserService.GetAccountInfo(c, &accountApiModel.GetUserInfoReq{AccountId: accountId})
	xgin.ResponseWithLang(c, resp, err)
}

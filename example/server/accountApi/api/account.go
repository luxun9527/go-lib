package api

import "github.com/gin-gonic/gin"

var UserApi = &userApi{}

type userApi struct{}

func (*userApi) AddAuth(c *gin.Context) {

}

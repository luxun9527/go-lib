package service

import (
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/model"
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/service/user"
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/service/user/impl"
)

var (
	UserRepository user.Repository
)
//Init instantiate the service
func Init()  {
	UserRepository = impl.NewMysqlImpl(model.MysqlHandler)
}
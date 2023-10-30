package userroleservicelogic

import (
	"context"

	"go-lib/frame/go-zero/rpcmult/internal/svc"
	"go-lib/frame/go-zero/rpcmult/mult"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleAddLogic {
	return &UserRoleAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleAddLogic) UserRoleAdd(in *mult.UserRoleAddReq) (*mult.UserRoleAddResp, error) {
	// todo: add your logic here and delete this line

	return &mult.UserRoleAddResp{}, nil
}

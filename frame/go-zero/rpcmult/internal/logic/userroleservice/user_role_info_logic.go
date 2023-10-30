package userroleservicelogic

import (
	"context"

	"go-lib/frame/go-zero/rpcmult/internal/svc"
	"go-lib/frame/go-zero/rpcmult/mult"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleInfoLogic {
	return &UserRoleInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleInfoLogic) UserRoleInfo(in *mult.UserRoleInfoReq) (*mult.UserRoleInfoResp, error) {
	// todo: add your logic here and delete this line

	return &mult.UserRoleInfoResp{}, nil
}

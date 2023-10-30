package userroleservicelogic

import (
	"context"

	"go-lib/frame/go-zero/rpcmult/internal/svc"
	"go-lib/frame/go-zero/rpcmult/mult"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleDeleteLogic {
	return &UserRoleDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleDeleteLogic) UserRoleDelete(in *mult.UserRoleDeleteReq) (*mult.UserRoleDeleteResp, error) {
	// todo: add your logic here and delete this line

	return &mult.UserRoleDeleteResp{}, nil
}

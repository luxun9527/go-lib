package userroleservicelogic

import (
	"context"

	"go-lib/frame/go-zero/rpcmult/internal/svc"
	"go-lib/frame/go-zero/rpcmult/mult"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleListLogic {
	return &UserRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleListLogic) UserRoleList(in *mult.UserRoleListReq) (*mult.UserRoleListResp, error) {
	// todo: add your logic here and delete this line

	return &mult.UserRoleListResp{}, nil
}

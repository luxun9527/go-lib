package userservicelogic

import (
	"context"

	"go-lib/frame/go-zero/rpcmult/internal/svc"
	"go-lib/frame/go-zero/rpcmult/mult"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoUpdateLogic {
	return &UserInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoUpdateLogic) UserInfoUpdate(in *mult.UserInfoUpdateReq) (*mult.UserInfoUpdateResp, error) {
	// todo: add your logic here and delete this line

	return &mult.UserInfoUpdateResp{}, nil
}

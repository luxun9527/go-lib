package logic

import (
	"context"

	"go-lib/frame/go-zero/apidemo/internal/svc"
	"go-lib/frame/go-zero/apidemo/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByNameLogic {
	return &GetUserInfoByNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoByNameLogic) GetUserInfoByName(req *types.GetUserInfoByNameReq) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}

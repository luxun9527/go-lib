package logic

import (
	"context"

	"go-lib/frame/go-zero/rpcdemo/internal/svc"
	"go-lib/frame/go-zero/rpcdemo/rpcdemo"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *rpcdemo.Request) (*rpcdemo.Response, error) {
	// todo: add your logic here and delete this line

	return &rpcdemo.Response{}, nil
}

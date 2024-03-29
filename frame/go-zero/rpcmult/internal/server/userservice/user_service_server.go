// Code generated by goctl. DO NOT EDIT.
// Source: mult.proto

package server

import (
	"context"

	"go-lib/frame/go-zero/rpcmult/internal/logic/userservice"
	"go-lib/frame/go-zero/rpcmult/internal/svc"
	"go-lib/frame/go-zero/rpcmult/mult"
)

type UserServiceServer struct {
	svcCtx *svc.ServiceContext
	mult.UnimplementedUserServiceServer
}

func NewUserServiceServer(svcCtx *svc.ServiceContext) *UserServiceServer {
	return &UserServiceServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServiceServer) Login(ctx context.Context, in *mult.LoginReq) (*mult.LoginResp, error) {
	l := userservicelogic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *UserServiceServer) UserInfo(ctx context.Context, in *mult.UserInfoReq) (*mult.UserInfoResp, error) {
	l := userservicelogic.NewUserInfoLogic(ctx, s.svcCtx)
	return l.UserInfo(in)
}

func (s *UserServiceServer) UserInfoUpdate(ctx context.Context, in *mult.UserInfoUpdateReq) (*mult.UserInfoUpdateResp, error) {
	l := userservicelogic.NewUserInfoUpdateLogic(ctx, s.svcCtx)
	return l.UserInfoUpdate(in)
}

func (s *UserServiceServer) UserList(ctx context.Context, in *mult.UserListReq) (*mult.UserListResp, error) {
	l := userservicelogic.NewUserListLogic(ctx, s.svcCtx)
	return l.UserList(in)
}

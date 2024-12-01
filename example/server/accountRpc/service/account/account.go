package account

import (
	"context"
	"errors"
	"github.com/gookit/goutil/strutil"
	"github.com/luxun9527/zlog"
	"go-lib/example/common/errs"
	accountPb "go-lib/example/pb/account"
	"go-lib/example/pkg/xjwt"
	"go-lib/example/server/accountRpc/common/model"
	"go-lib/example/server/accountRpc/common/rediskey"
	"go-lib/example/server/accountRpc/global"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"
)

var AccountRpc accountRpc

type accountRpc struct {
	accountPb.UnimplementedAccountSrvServer
}

func (accountRpc) GetAllUserInfo(context.Context, *emptypb.Empty) (*accountPb.GetAccountInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllUserInfo not implemented")
}
func (accountRpc) GetAccountInfo(ctx context.Context, req *accountPb.GetAccountInfoReq) (*accountPb.GetAccountInfoResp, error) {
	account := global.AccountDB.Account
	accountInfo, err := account.WithContext(ctx).Where(account.AccountID.Eq(req.AccountId)).Take()
	if err != nil {
		zlog.Errorf("GetAccountInfo err: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.RecordNotFoundErr
		}
		return nil, err
	}
	return &accountPb.GetAccountInfoResp{
		AccountId:   accountInfo.AccountID,
		AccountName: accountInfo.AccountName,
	}, nil
}
func (accountRpc) RegisterUser(context.Context, *accountPb.RegisterUserReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (accountRpc) ValidateToken(ctx context.Context, req *accountPb.ValidateTokenReq) (*accountPb.ValidateTokenResp, error) {
	userInfo, err := xjwt.ParseToken[model.UserInfo](req.Token)
	if err != nil {
		zlog.Errorf("ValidateToken err: %v", err)
		return nil, err
	}
	tokenMd5 := strutil.Md5(req.Token)
	c, err := global.RedisCli.Exists(ctx, rediskey.AccountToken.WithParams(tokenMd5)).Result()
	if err != nil {
		zlog.Errorf("ValidateToken exist redis key err: %v", err)
		return nil, errs.RedisErr
	}
	if c == 0 {
		return nil, errs.TokenValidateFailed
	}

	return &accountPb.ValidateTokenResp{
		AccountId:   userInfo.Extra.AccountId,
		AccountName: userInfo.Extra.AccountName,
	}, nil
}
func (accountRpc) Login(ctx context.Context, req *accountPb.LoginReq) (*accountPb.LoginResp, error) {
	account := global.AccountDB.Account
	req.Password = strutil.Md5(req.Password)
	accountInfo, err := account.WithContext(ctx).Where(account.AccountID.Eq(req.AccountName), account.Password.Eq(req.Password)).Take()
	if err != nil {
		zlog.Errorf("Login err: %v", err)
		return nil, errs.LoginFailed
	}
	userInfo := model.UserInfo{
		AccountId:   accountInfo.AccountID,
		AccountName: accountInfo.AccountName,
	}
	c, err := xjwt.NewCustomClaims[model.UserInfo](userInfo)
	if err != nil {
		return nil, err
	}
	token, err := c.GenerateToken()
	if err != nil {
		return nil, err
	}

	tokenMd5 := strutil.Md5(token)
	if err := global.RedisCli.Set(ctx, rediskey.AccountToken.WithParams(tokenMd5), "", time.Hour*24).Err(); err != nil {
		zlog.Errorf("Login set redis key err: %v", err)
		return nil, errs.RedisErr
	}
	return &accountPb.LoginResp{
		Token: token,
	}, nil
}

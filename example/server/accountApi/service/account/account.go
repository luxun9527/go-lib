package account

import (
	"context"
	"github.com/luxun9527/zlog"
	accountPb "go-lib/example/pb/account"
	accountApiModel "go-lib/example/server/accountApi/model/account"
	"go-lib/example/server/accountApi/rpcClient"
)

var UserService = &userService{}

type userService struct{}

func (*userService) GetAccountInfo(ctx context.Context, req *accountApiModel.GetUserInfoReq) (*accountApiModel.GetUserInfoResp, error) {
	accountInfo, err := rpcClient.AccountClient.GetAccountInfo(ctx, &accountPb.GetAccountInfoReq{AccountId: req.AccountId})
	if err != nil {
		zlog.Errorf("GetAccountInfo err: %v", err)
		return nil, err
	}
	return &accountApiModel.GetUserInfoResp{
		AccountId:   accountInfo.AccountId,
		AccountName: accountInfo.AccountName,
	}, nil
}

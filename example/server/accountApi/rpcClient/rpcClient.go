package rpcClient

import (
	"context"
	accountPb "go-lib/example/pb/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	AccountClient accountPb.AccountSrvClient
)

func InitRpcClient() error {
	// 初始化RPC客户端
	InitAccountRpcClient(context.Background(), "localhost:50051")
	return nil
}
func InitAccountRpcClient(ctx context.Context, addr string) error {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	AccountClient = accountPb.NewAccountSrvClient(conn)
	return nil
}

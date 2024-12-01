package rpcClient

import (
	"context"
	"github.com/luxun9527/zlog"
	accountPb "go-lib/example/pb/account"
	"go-lib/example/pkg/grpcx/clientinterceptors"
	"go-lib/example/server/accountApi/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	AccountClient accountPb.AccountSrvClient
)
var (
	InitClientFuncs = []func(conn *grpc.ClientConn){
		func(conn *grpc.ClientConn) {
			AccountClient = accountPb.NewAccountSrvClient(conn)
		},
	}
)

func InitEtcdRpcClients(ctx context.Context, cli *clientv3.Client, targetConfList []*config.TargetConf) error {
	// 初始化RPC客户端
	for i, v := range targetConfList {
		grpcConn, err := initEtcdGrpcConn(ctx, cli, "etcd:///"+v.Key, v.TimeOut)
		if err != nil {
			return err
		}
		InitClientFuncs[i](grpcConn)
	}
	return nil
}

func initEtcdGrpcConn(ctx context.Context, cli *clientv3.Client, addr string, timeOut int64) (*grpc.ClientConn, error) {
	etcdResolver, err := resolver.NewBuilder(cli)
	if err != nil {
		zlog.Panicf("create etcd resolver failed, err:%v", err)
	}
	return grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithChainUnaryInterceptor(clientinterceptors.TimeoutInterceptor(time.Second*time.Duration(timeOut)), clientinterceptors.LoggerInterceptor()),
	)

}

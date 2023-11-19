package main

import (
	"flag"
	"fmt"
	"go-lib/frame/go-zero/rpcmult/internal/config"
	userroleserviceServer "go-lib/frame/go-zero/rpcmult/internal/server/userroleservice"
	userserviceServer "go-lib/frame/go-zero/rpcmult/internal/server/userservice"
	"go-lib/frame/go-zero/rpcmult/internal/svc"
	"go-lib/frame/go-zero/rpcmult/mult"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/mult.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		mult.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))
		mult.RegisterUserRoleServiceServer(grpcServer, userroleserviceServer.NewUserRoleServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

package main

import (
	"flag"
	"fmt"
	"github.com/luxun9527/zlog"
	accountPb "go-lib/example/pb/account"
	"go-lib/example/server/accountRpc/global"
	"go-lib/example/server/accountRpc/initializer"
	accountService "go-lib/example/server/accountRpc/service/account"
	"google.golang.org/grpc"
	"net"
)

var (
	path = flag.String("f", "example/server/accountRpc/conf/config.yaml", "config file path")
)

func main() {
	flag.Parse()
	initializer.Init(*path)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", global.Config.Server.Port))
	if err != nil {
		zlog.Panicf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	accountPb.RegisterAccountSrvServer(s, accountService.AccountRpc)
	zlog.Infof("start rpc server on %d", global.Config.Server.Port)
	if err := s.Serve(listener); err != nil {
		zlog.Panicf("failed to serve: %v", err)
	}
}

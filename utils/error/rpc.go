package error

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"
)

type DemoServer struct {
	UnimplementedErrorsServer
}

func (DemoServer) Demo(ctx context.Context, req *Empty) (*Empty, error) {
	val := time.Now().UnixMilli() % 3
	switch val {
	case 2:
		//模拟异常，返回未知错误，返回给api,api解析为code.Unknown 2
		return nil, errors.New("unknown error")
	case 1:
		//模拟异常，grpc中的code。返回给api,api解析为5
		return nil, status.Error(codes.NotFound, "user not found")
	default:
		//业务异常，返回自定义错误
		return nil, UserNotFound
	}

}
func InitRpcSrv() {
	listener, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		log.Panicf("failed to start listener: %v", err)
	}
	s := grpc.NewServer()
	RegisterErrorsServer(s, &DemoServer{})
	log.Printf("start rpc server on 0.0.0.0:8899")
	if err := s.Serve(listener); err != nil {
		log.Panicf("failed to start rpc server: %v", err)
	}

}

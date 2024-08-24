package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-lib/net/grpc/pb/grpcdemo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	//开启反射
	reflection.Register(s)
	grpcdemo.RegisterGrpcDemoServer(s, new(GrpcDemoServer))
	log.Printf("start server at %v", 8899)
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}

}

func TestGrpcGateWayServer(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:8899")
		if err != nil {
			log.Println("net listen err ", err)
			return
		}
		s := grpc.NewServer()
		grpcdemo.RegisterGrpcDemoServer(s, new(GrpcDemoServer))
		grpcdemo.RegisterGrpcGatewayDemoServer(s, new(GrpcGatewayDemo))
		if err := s.Serve(listener); err != nil {
			log.Println("failed to serve...", err)
			return
		}
	}()
	conn, err := grpc.Dial(
		"127.0.0.1:8899",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Panic("dail proxy grpc serve failed ", zap.Error(err))
	}

	gwmux := runtime.NewServeMux()

	if err = grpcdemo.RegisterGrpcDemoHandler(context.Background(), gwmux, conn); err != nil {
		log.Panicf("Failed to register gateway %v", err)
	}
	if err = grpcdemo.RegisterGrpcGatewayDemoHandler(context.Background(), gwmux, conn); err != nil {
		log.Panicf("Failed to register gateway %v", err)
	}

	gwServer := &http.Server{
		Addr:    ":10080",
		Handler: gwmux,
	}
	if err := gwServer.ListenAndServe(); err != nil {
		log.Panic("init proxy http serve failed err", zap.Error(err))
	}

}

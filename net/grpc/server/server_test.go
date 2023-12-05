package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
)

type GrpcDemoServer struct {
	port int32
	grpcdemo.UnimplementedGrpcDemoServer
}

func (c GrpcDemoServer) UnaryCall(ctx context.Context, req *emptypb.Empty) (*grpcdemo.UnaryCallResp, error) {
	log.Printf("port is %v", c.port)
	return &grpcdemo.UnaryCallResp{Username: "zhangsan"}, nil
}
func (c GrpcDemoServer) DemoImport(ctx context.Context, req *folder.ImportedMessage) (*grpcdemo.CustomMessage, error) {
	log.Printf("port is %v", c.port)
	return &grpcdemo.CustomMessage{
		CustomMessage: fmt.Sprintf("%v", c.port),
	}, nil
}
func (GrpcDemoServer) PushData(c grpcdemo.GrpcDemo_PushDataServer) error {
	for {
		data, err := c.Recv()
		if err != nil {
			log.Printf("err %v", err)
			return err
		}
		log.Printf("recv data %v", data)
	}

}
func (GrpcDemoServer) FetchData(req *grpcdemo.FetchDataReq, c grpcdemo.GrpcDemo_FetchDataServer) error {
	for i := 0; i < 10; i++ {
		if err := c.Send(&grpcdemo.FetchDataResp{
			FavBook: "book",
		}); err != nil {
			log.Printf("err %v", err)
			return err
		}
	}
	return nil
}
func (GrpcDemoServer) Exchange(c grpcdemo.GrpcDemo_ExchangeServer) error {
	g := sync.WaitGroup{}

	g.Add(2)
	go func() {
		defer g.Done()
		for {
			data, err := c.Recv()
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("exchange recv message %v", data)
		}

	}()
	go func() {
		defer g.Done()
		for {
			if err := c.Send(&grpcdemo.ExchangeResp{LastName: "test"}); err != nil {
				if err != nil {
					log.Println(err)
					return
				}
			}
		}

	}()
	g.Wait()
	return nil
}

func (GrpcDemoServer) CallGrpcGateway(ctx context.Context, req *grpcdemo.CallGrpcGatewayReq) (*grpcdemo.CallGrpcGatewayResp, error) {
	log.Printf("recv message %v", req.Config)
	name := req.Config["name"]
	switch name {
	case "zhangsan":
		return nil, status.Error(codes.NotFound, "not found")
	case "lisi":
		return nil, errors.New("this is custom error")

	}
	return &grpcdemo.CallGrpcGatewayResp{Config: req.Config}, nil
}

type GrpcGatewayDemo struct {
	grpcdemo.GrpcGatewayDemoServer
}

func (GrpcGatewayDemo) CallGrpcGatewayDemo(ctx context.Context, req *grpcdemo.CallGrpcGatewayDemoReq) (*grpcdemo.CallGrpcGatewayDemoResp, error) {

	return &grpcdemo.CallGrpcGatewayDemoResp{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}
	s := grpc.NewServer()
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

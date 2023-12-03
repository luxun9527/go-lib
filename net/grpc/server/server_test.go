package server

import (
	"context"
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

func (c GrpcDemoServer) UnaryCall(ctx context.Context, req *grpcdemo.NoticeReaderReq) (*emptypb.Empty, error) {
	log.Printf("port is %v", c.port)
	return &emptypb.Empty{}, nil
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
		log.Println(data)
	}

}
func (GrpcDemoServer) FetchData(req *grpcdemo.Empty, c grpcdemo.GrpcDemo_FetchDataServer) error {
	for i := 0; i < 10; i++ {
		if err := c.Send(&grpcdemo.Data{}); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
func (GrpcDemoServer) Exchange(c grpcdemo.GrpcDemo_ExchangeServer) error {
	g := sync.WaitGroup{}

	g.Add(2)
	go func() {
		for {
			req, err := c.Recv()
			if err != nil {
				log.Println(err)
			}
			log.Println(req)
		}
		g.Done()

	}()
	go func() {
		for {
			if err := c.Send(&grpcdemo.Resp{LastName: "test"}); err != nil {
				if err != nil {
					log.Println(err)
				}
			}
		}
		g.Done()
	}()
	g.Wait()
	return nil
}

func (GrpcDemoServer) CallGrpcGateway(ctx context.Context, req *grpcdemo.NoticeReaderReq) (*grpcdemo.NoticeReaderResp, error) {
	log.Println(req)
	switch req.Msg {
	case "1":
		return nil, status.Error(codes.NotFound, "not found")
	case "2":
		return nil, fmt.Errorf("custom error")

	}
	return &grpcdemo.NoticeReaderResp{FavBook: ""}, nil
}

type GrpcGatewayDemo struct {
	grpcdemo.GrpcGatewayDemoServer
}

func (GrpcGatewayDemo) CallGrpcGatewayDemo(ctx context.Context, req *grpcdemo.NoticeReaderReq) (*grpcdemo.NoticeReaderResp, error) {
	switch req.Msg {
	case "1":
		return nil, status.Error(codes.NotFound, "CallGrpcGatewayDemo not found")
	case "2":
		return nil, fmt.Errorf("CallGrpcGatewayDemo custom error")

	}
	return &grpcdemo.NoticeReaderResp{FavBook: "test11"}, nil
}

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}
	s := grpc.NewServer()
	grpcdemo.RegisterGrpcDemoServer(s, new(GrpcDemoServer))
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

	go func() {
		if err := gwServer.ListenAndServe(); err != nil {
			log.Panic("init proxy http serve failed err", zap.Error(err))
		}
	}()
	select {}
}

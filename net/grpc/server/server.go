package server

import (
	"context"
	"errors"
	"fmt"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"sync"
)

type GrpcDemoServer struct {
	port int32
	grpcdemo.UnimplementedGrpcDemoServer
}

func (c GrpcDemoServer) UnaryCall(ctx context.Context, req *emptypb.Empty) (*grpcdemo.UnaryCallResp, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get interface addresses: %v", err)
	}
	ips := ""
	for _, addr := range addrs {
		// 检查地址类型并跳过回环地址
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// 仅打印IPv4地址
			if ipNet.IP.To4() != nil {

				ips += "_" + ipNet.IP.String()
			}
		}
	}
	return &grpcdemo.UnaryCallResp{Username: "zhangsan,ips:" + ips}, nil
}
func (c GrpcDemoServer) DemoImport(ctx context.Context, req *folder.ImportedMessage) (*grpcdemo.CustomMessage, error) {
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
		Data: &anypb.Any{
			TypeUrl: "",
			Value:   nil,
		},
	}, nil
}

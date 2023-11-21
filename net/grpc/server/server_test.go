package server

import (
	"context"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"sync"
	"testing"
)

type GrpcDemoServer struct {
	grpcdemo.UnimplementedGrpcDemoServer
}
func (GrpcDemoServer) Call(ctx context.Context,req *grpcdemo.NoticeReaderReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (GrpcDemoServer) DemoImport(ctx context.Context, req *folder.ImportedMessage) (*grpcdemo.CustomMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DemoImport not implemented")
}
func (GrpcDemoServer) PushData(c grpcdemo.GrpcDemo_PushDataServer) error {
	for  {
		data, err := c.Recv()
		if err!=nil{
			log.Printf("err %v",err)
			return err
		}
		log.Println(data)
	}


}
func (GrpcDemoServer) FetchData(req *grpcdemo.Empty,c grpcdemo.GrpcDemo_FetchDataServer) error {
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
		for  {
			req, err := c.Recv()
			if err != nil{
				log.Println(err)
			}
			log.Println(req)
		}
		g.Done()

	}()
	go func() {
		for  {
			if err := c.Send(&grpcdemo.Resp{LastName: "test"});err!=nil{
				if err != nil{
					log.Println(err)
				}
			}
		}
		g.Done()
	}()
	g.Wait()
	return nil
}

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}
	s := grpc.NewServer()
	grpcdemo.RegisterGrpcDemoServer(s,new(GrpcDemoServer))
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}
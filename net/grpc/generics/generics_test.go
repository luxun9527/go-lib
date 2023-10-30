package generics

import (
	"context"
	hellopb "go-lib/net/grpc/hello/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Hello struct {
	hellopb.UnimplementedHelloServiceServer
}

func (Hello) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {

	log.Printf("receive message %+v\n", req)
	//return nil, NotFound
	return &hellopb.HelloResponse{
		Reply:      "hello",
		SecondName: "zhangsan",
	}, nil
}

func RunServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	hellopb.RegisterHelloServiceServer(s, new(Hello))
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}

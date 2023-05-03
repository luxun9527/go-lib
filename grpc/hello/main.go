package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	hellopb "hello/pb"
	"log"
	"net"
	"time"
)

func main() {
	go RunServer()
	time.Sleep(time.Second * 1)
	RunClient()
}

type Hello struct {
	hellopb.UnimplementedHelloServiceServer
}

func (Hello) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	fmt.Printf("receive message %+v\n", req)

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
func RunClient() {
	conn, err := grpc.Dial("0.0.0.0:9090", grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect.", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		ctx := context.Background()
		helloClient := hellopb.NewHelloServiceClient(conn)
		resp, err := helloClient.SayHello(ctx, &hellopb.HelloRequest{
			Name: "zhangsan",
			Age:  12,
		})
		if err != nil {
			fmt.Printf("error = %v", err)
			return
		}
		fmt.Printf("resp = %+v\n", resp)

	}

}

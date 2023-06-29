package main

import (
	"context"
	hellopb "go-lib/net/grpc/hello/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"log"
	"net"
	"sync"
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

	log.Printf("receive message %+v\n", req)
	return nil, status.New(codes.Code(100000), "").Err()
	//return &hellopb.HelloResponse{
	//	Reply:      "hello",
	//	SecondName: "zhangsan",
	//}, nil
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
	helloClient := hellopb.NewHelloServiceClient(conn)
	group := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		group.Add(i)
		go func() {
			defer group.Done()
			ctx := context.Background()
			resp, err := helloClient.SayHello(ctx, &hellopb.HelloRequest{
				Name: "zhangsan",
				Age:  12,
			})
			if err != nil {
				log.Println(err)

			}
			log.Printf("resp = %+v\n", resp)
		}()

	}
	group.Wait()

}

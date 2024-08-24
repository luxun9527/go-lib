package main

import (
	"context"
	"flag"
	"github.com/sercand/kuberesolver/v5"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

var (
	mode = flag.String("mode", "server", "server or client")
)

func main() {
	flag.Parse()
	log.Println("mode:", *mode)
	if *mode == "client" {
		initClient()
	} else {
		initServer()
	}

}
func initServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	//开启反射
	reflection.Register(s)
	grpcdemo.RegisterGrpcDemoServer(s, new(server.GrpcDemoServer))
	log.Printf("start server at %v", 8899)
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}
func initClient() {
	kuberesolver.RegisterInCluster()

	conn, err := grpc.Dial("kubernetes:///grpcdemoserver-service.default:8899",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)

	for {
		time.Sleep(time.Second * 5)
		result, err := cli.UnaryCall(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Printf("Call  failed %v", err)
		} else {
			log.Printf("resp %v", result)
		}
	}
}

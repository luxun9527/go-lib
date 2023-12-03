package server

import (
	"go-lib/net/grpc/pb/grpcdemo"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"testing"
)

func TestResolverServer(t *testing.T) {
	w := sync.WaitGroup{}
	w.Add(3)
	go func() {
		defer w.Done()
		listener, err := net.Listen("tcp", "0.0.0.0:8898")
		if err != nil {
			log.Println("net listen err ", err)
			return
		}
		s := grpc.NewServer()
		grpcdemo.RegisterGrpcDemoServer(s, &GrpcDemoServer{
			port:                        8898,
			UnimplementedGrpcDemoServer: grpcdemo.UnimplementedGrpcDemoServer{},
		})
		log.Printf("start server at %v", 8898)
		if err := s.Serve(listener); err != nil {
			log.Println("failed to serve...", err)
			return
		}
	}()

	go func() {
		defer w.Done()
		listener, err := net.Listen("tcp", "0.0.0.0:8899")
		if err != nil {
			log.Println("net listen err ", err)
			return
		}
		s := grpc.NewServer()
		grpcdemo.RegisterGrpcDemoServer(s, &GrpcDemoServer{
			port:                        8899,
			UnimplementedGrpcDemoServer: grpcdemo.UnimplementedGrpcDemoServer{},
		})
		log.Printf("start server at %v", 8899)
		if err := s.Serve(listener); err != nil {
			log.Println("failed to serve...", err)
			return
		}

	}()
	go func() {
		defer w.Done()
		listener, err := net.Listen("tcp", "0.0.0.0:8897")
		if err != nil {
			log.Println("net listen err ", err)
			return
		}
		s := grpc.NewServer()
		grpcdemo.RegisterGrpcDemoServer(s, &GrpcDemoServer{
			port:                        8897,
			UnimplementedGrpcDemoServer: grpcdemo.UnimplementedGrpcDemoServer{},
		})
		log.Printf("start server at %v", 8897)
		if err := s.Serve(listener); err != nil {
			log.Println("failed to serve...", err)
			return
		}

	}()
	w.Wait()
}

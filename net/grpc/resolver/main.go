package main

import (
	"context"
	"go-lib/net/grpc/resolver/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
)

// 自定义name resolver

const (
	myScheme   = "q1mi"
	myEndpoint = "resolver.liwenzhou.com"
)

var addrs = []string{"127.0.0.1:8972", "127.0.0.1:8973"}

// q1miResolver 自定义name resolver，实现Resolver接口
type q1miResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *q1miResolver) ResolveNow(o resolver.ResolveNowOptions) {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*q1miResolver) Close() {}

// q1miResolverBuilder 需实现 Builder 接口
type q1miResolverBuilder struct{}

func (*q1miResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &q1miResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			myEndpoint: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*q1miResolverBuilder) Scheme() string { return myScheme }

func main() {
	//	go InitServer1()
	//go InitServer2()

	//for i := 0; i < 10; i++ {
	//	conn, err := grpc.Dial(
	//		"q1mi:///resolver.liwenzhou.com",
	//		grpc.WithTransportCredentials(insecure.NewCredentials()),
	//		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 这里设置初始策略
	//		grpc.WithResolvers(&q1miResolverBuilder{}),                             // 指定使用q1miResolverBuilder
	//	)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	c := pb.NewResolverServiceClient(conn)
	//	resp, err := c.Resolver(context.Background(), &pb.ResolverReq{
	//		Name: "",
	//		Age:  0,
	//	})
	//	if err != nil {
	//		log.Println(err)
	//	} else {
	//		log.Println(resp)
	//	}
	//}
	type name struct {
		age string
	}
	m := make(map[name]string)
	m[name{age: "test"}] = ""
	for k, v := range m {
		log.Println(k, v)
	}
	log.Println(m[name{age: "test"}])
	//select {}
}
func InitServer1() {
	listener, err := net.Listen("tcp", "0.0.0.0:8973")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterResolverServiceServer(s, Resolver{ip: "8973"})
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}
func InitServer2() {
	listener, err := net.Listen("tcp", "0.0.0.0:8972")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterResolverServiceServer(s, Resolver{ip: "8972"})
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}

type Resolver struct {
	ip string
	pb.UnimplementedResolverServiceServer
}

func (r Resolver) Resolver(ctx context.Context, req *pb.ResolverReq) (*pb.ResolverResp, error) {
	log.Println("invoke", r.ip)
	return &pb.ResolverResp{
		Reply:      "11",
		SecondName: "11",
	}, nil
}

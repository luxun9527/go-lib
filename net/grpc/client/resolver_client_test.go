package client

import (
	"context"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"testing"
	"time"
)

/*
refer
https://www.liwenzhou.com/posts/Go/name-resolving-and-load-balancing-in-grpc/
https://github.com/win5do/go-microservice-demo/blob/main/docs/sections/grpc-lb.md
https://www.cnblogs.com/FireworksEasyCool/p/12912839.html
*/
// 自定义name resolver

var (
	_customScheme   = "dns"
	_customEndpoint = "xxx.xxx.com"
	_addrs          = []string{"127.0.0.1:8898", "127.0.0.1:8899"}
)

// customResolver 自定义name resolver，实现Resolver接口
type customResolver struct {
}

func (r *customResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*customResolver) Close() {}

// customBuilder 需实现 Builder 接口
type customBuilder struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (builder *customBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	builder.target = target
	builder.cc = cc
	builder.addrsStore = map[string][]string{_customEndpoint: _addrs}
	addresses := builder.addrsStore[target.Endpoint()]
	r := make([]resolver.Address, 0, 2)
	for _, v := range addresses {
		a := resolver.Address{
			Addr:       v,
			ServerName: "",
		}
		r = append(r, a)
	}
	if err := builder.cc.UpdateState(resolver.State{Addresses: r}); err != nil {
		return nil, err
	}

	go func() {
		time.Sleep(time.Second * 20)
		builder.updateConn()
	}()
	return &customResolver{}, nil
}
func (*customBuilder) Scheme() string { return _customScheme }

// 执行UpdateState更新连接
func (builder *customBuilder) updateConn() {
	addresses := []resolver.Address{{
		Addr:               "127.0.0.1:8898",
		ServerName:         "",
		Attributes:         nil,
		BalancerAttributes: nil,
		Metadata:           nil,
	}, {
		Addr:               "127.0.0.1:8899",
		ServerName:         "",
		Attributes:         nil,
		BalancerAttributes: nil,
		Metadata:           nil,
	}, {
		Addr:               "127.0.0.1:8897",
		ServerName:         "",
		Attributes:         nil,
		BalancerAttributes: nil,
		Metadata:           nil,
	}}
	if err := builder.cc.UpdateState(resolver.State{Addresses: addresses}); err != nil {
		log.Printf("update state failed %v", err)
	}
}

func TestResolverClientTest(t *testing.T) {
	resolverBuilder := &customBuilder{}
	conn, err := grpc.Dial(
		"dns:///xxx.xxx.com",
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(resolverBuilder),
	)
	if err != nil {
		log.Printf("dial connection failed err =%v", err)
	}

	cli := grpcdemo.NewGrpcDemoClient(conn)
	for i := 0; i < 200; i++ {
		resp, err := cli.DemoImport(context.Background(), &folder.ImportedMessage{
			ImportedMessage: "test",
		})
		if err != nil {
			log.Printf("call failed err %v", err)
		} else {
			log.Printf("resp %v", resp)
		}
		time.Sleep(time.Second)
	}

}

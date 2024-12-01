package test

import (
	"context"
	"github.com/zeromicro/go-zero/core/netx"
	"github.com/zeromicro/go-zero/core/stringx"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestRegister1(t *testing.T) {
	e := []string{"192.168.2.159:2379"}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Panicf("connect to etcd failed, err:%v", err)
	}
	manager, err := endpoints.NewManager(cli, "xxRpc")
	if err != nil {
		log.Panicf("create service failed, err:%v", err)
	}
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		log.Panicf("create lease failed, err:%v", err)
	}
	if err := manager.AddEndpoint(context.Background(), "xxRpc/"+stringx.Randn(10), endpoints.Endpoint{
		Addr:     netx.InternalIp() + ":8899",
		Metadata: nil,
	}, clientv3.WithLease(resp.ID)); err != nil {
		log.Panicf("add endpoint failed, err:%v", err)
	}
	if err := manager.AddEndpoint(context.Background(), "xxRpc/"+stringx.Randn(10), endpoints.Endpoint{
		Addr:     netx.InternalIp() + ":8897",
		Metadata: nil,
	}, clientv3.WithLease(resp.ID)); err != nil {
		log.Panicf("add endpoint failed, err:%v", err)
	}
	c, err := cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		log.Panicf("keep alive failed, err:%v", err)
	}
	for _ = range c {

	}
	select {}
}
func TestResolver(t *testing.T) {
	e := []string{"192.168.2.159:2379"}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Panicf("connect to etcd failed, err:%v", err)
	}
	etcdResolver, err := resolver.NewBuilder(cli)
	if err != nil {
		log.Panicf("create etcd resolver failed, err:%v", err)
	}
	conn, err := grpc.Dial("etcd:///"+"xxRpc",
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithInsecure())
	if err != nil {
		log.Panicf("connect to etcd failed, err:%v", err)
	}
	client := grpcdemo.NewGrpcDemoClient(conn)
	for i := 0; i < 200; i++ {
		resp, err := client.DemoImport(context.Background(), &folder.ImportedMessage{
			ImportedMessage: "test",
		})
		if err != nil {
			log.Printf("call failed err %v", err)
		} else {
			log.Printf("resp %v", resp)
			//2024/11/27 23:02:03 resp custom_message:"8897"
			//2024/11/27 23:02:04 resp custom_message:"8898"
			//2024/11/27 23:02:05 resp custom_message:"8899"
			//2024/11/27 23:02:06 resp custom_message:"8897"
			//2024/11/27 23:02:07 resp custom_message:"8898"
			//2024/11/27 23:02:08 resp custom_message:"8899"
		}
		time.Sleep(time.Second)
	}
}

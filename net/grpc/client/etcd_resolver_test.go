package client

import (
	"context"
	ws "github.com/luxun9527/gpush/proto"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

var (
	service = "resolver"
	addr    = "192.168.2.138:8973"
)

func TestRegister(t *testing.T) {
	if err := register(); err != nil {
		log.Println(err)
		return
	}
	select {}
}

func TestResolve(t *testing.T) {
	resolve()

}
func resolve() {
	e := []string{"192.168.2.99:2379"}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	etcdResolver, err := NewBuilder(cli)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.Dial("etcd:///"+"proxy", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithResolvers(etcdResolver))
	if err != nil {
		log.Println(err)
	}

	client := ws.NewProxyClient(conn)
	if _, err := client.PushData(context.Background(), &ws.Data{
		Uid:   "",
		Topic: "hello",
		Data:  []byte("abcd"),
	}); err != nil {
		log.Println(err)
	}
}
func register() error {
	e := []string{"192.168.2.99:2379"}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e,
		DialTimeout: 5 * time.Second,
	})
	//设置租约时间
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = cli.Put(context.Background(), service, addr, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	keepAliveresp, err := cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		log.Println(err)
	}
	for _ = range keepAliveresp {

	}
	return nil

}

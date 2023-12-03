package client

import (
	"context"
	"github.com/spf13/cast"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	gresolver "google.golang.org/grpc/resolver"
	"log"
	"sync"
	"testing"
	"time"
)

type etcdBuilder struct {
	etcdCli    *clientv3.Client
	cc         gresolver.ClientConn
	target     gresolver.Target
	serverList map[string]gresolver.Address
	lock       sync.Mutex
}

func (builder *etcdBuilder) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (gresolver.Resolver, error) {
	endpoint := target.Endpoint()
	builder.target = target
	builder.cc = cc
	builder.serverList = map[string]gresolver.Address{}

	resp, err := builder.etcdCli.Get(context.Background(), endpoint, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	addrs := make([]gresolver.Address, 0, len(resp.Kvs))
	for _, v := range resp.Kvs {
		a := gresolver.Address{
			Addr: string(v.Value),
		}
		addrs = append(addrs, a)
		builder.serverList[string(v.Key)] = a
	}
	if err := builder.cc.UpdateState(gresolver.State{Addresses: addrs}); err != nil {
		return nil, err
	}

	go builder.Watch(endpoint)
	return &etcdResolver{}, nil
}
func (builder *etcdBuilder) Watch(prefix string) {
	events := builder.etcdCli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for r := range events {
		for _, ev := range r.Events {
			switch ev.Type {
			case mvccpb.PUT:
				builder.lock.Lock()
				builder.serverList[string(ev.Kv.Key)] = gresolver.Address{Addr: string(ev.Kv.Value)}
				addrs := make([]gresolver.Address, 0, len(builder.serverList))
				for _, v := range builder.serverList {
					addrs = append(addrs, v)
				}
				if err := builder.cc.UpdateState(gresolver.State{Addresses: addrs}); err != nil {
					log.Printf("watch UpdateState failed %v", err)
				}
				builder.lock.Unlock()
			case mvccpb.DELETE:
				builder.lock.Lock()
				delete(builder.serverList, string(ev.Kv.Key))
				addrs := make([]gresolver.Address, 0, len(builder.serverList))
				for _, v := range builder.serverList {
					addrs = append(addrs, v)
				}
				if err := builder.cc.UpdateState(gresolver.State{Addresses: addrs}); err != nil {
					log.Printf("watch UpdateState failed %v", err)
				}
				builder.lock.Unlock()
			}
		}
	}
}
func removeElem(addrs []gresolver.Address, elem string) []gresolver.Address {
	j := 0
	for _, v := range addrs {
		if v.Addr != elem {
			addrs[j] = v
			j++
		}
	}
	return addrs[:j]
}

func (builder *etcdBuilder) Scheme() string {
	return "etcd"
}

// NewBuilder creates a etcdResolver etcdBuilder.
func NewBuilder(client *clientv3.Client) (gresolver.Builder, error) {
	return &etcdBuilder{etcdCli: client}, nil
}

type etcdResolver struct {
}

// ResolveNow is a no-op here.
// It's just a hint, etcdResolver can ignore this if it's not necessary.
func (r *etcdResolver) ResolveNow(gresolver.ResolveNowOptions) {}

func (r *etcdResolver) Close() {}

var (
	key = "etcdResolver"
)

func TestEtcdResolve(t *testing.T) {
	go func() {
		go func() {
			register("127.0.0.1:8898")
		}()
		time.Sleep(time.Second * 10)
		go func() {
			register("127.0.0.1:8899")
		}()
	}()

	time.Sleep(time.Second * 3)
	e := []string{"192.168.2.159:2379"}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Panicf("init client failed %v", err)
	}
	builder, err := NewBuilder(cli)
	if err != nil {
		log.Panicf("init build failed %v", err)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	conn, err := grpc.DialContext(ctx, "etcd:///"+key,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithResolvers(builder),
		grpc.WithBlock())
	if err != nil {
		log.Panicf("connect failed %v", err)
	}
	c := grpcdemo.NewGrpcDemoClient(conn)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		resp, err := c.DemoImport(context.Background(), &folder.ImportedMessage{})
		if err != nil {
			log.Println("unary call failed", err)
		} else {
			log.Printf("resp = %v", resp)
		}
	}

}

func register(addr string) {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.2.159:2379"},
		DialTimeout: 5 * time.Second,
	})
	//设置租约时间
	resp, err := cli.Grant(context.Background(), 5)
	if err != nil {
		log.Panicf("init client failed %v", err)
	}
	//注册服务并绑定租约
	_, err = cli.Put(context.Background(), key+"/"+cast.ToString(int64(resp.ID)), addr, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Panicf("register client failed %v", err)
	}

	c, err := cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		log.Panicf("KeepAlive  failed %v", err)
	}
	for _ = range c {

	}

}

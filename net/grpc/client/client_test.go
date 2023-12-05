package client

import (
	"context"
	"go-lib/net/grpc/pb/grpcdemo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)

	for {
		time.Sleep(time.Second * 10)
		result, err := cli.UnaryCall(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Printf("Call  failed %v", err)
		} else {
			log.Printf("resp %v", result)
		}
	}
}

func TestPush(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)

	c, err := cli.PushData(context.Background())
	if err != nil {
		log.Printf("get pushdata connection failed %v", err)
		return
	}
	for i := 0; i < 10; i++ {
		if err := c.Send(&grpcdemo.PushDataReq{Foo: "foo"}); err != nil {
			log.Printf("push data failed %v", err)
			return
		}
	}

}

func TestFetchData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)
	c, err := cli.FetchData(context.Background(), &grpcdemo.FetchDataReq{
		Msg:       "",
		NoticeWay: &grpcdemo.FetchDataReq_Email{Email: "test"},
	})
	if err != nil {
		log.Printf("get fetchdata connection failed %v", err)
		return
	}

	for {
		data, err := c.Recv()
		if err != nil {
			log.Printf("recv data failed %v", err)
			return
		}
		log.Printf("data =%v", data)
	}

}
func TestExchangeData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)
	c, err := cli.Exchange(context.Background())
	if err != nil {
		log.Printf("get Exchangedata connection failed %v", err)
		return
	}
	group := sync.WaitGroup{}
	group.Add(2)
	go func() {
		for {
			defer group.Done()
			data, err := c.Recv()
			if err != nil {
				log.Printf("recv data failed %v", err)
				return
			}
			log.Printf("data =%v", data)
		}
	}()
	go func() {
		for {
			defer group.Done()
			var age = "12"
			err := c.Send(&grpcdemo.ExchangeReq{
				FirstName: "zhangsan",
				Age:       &age,
			})
			if err != nil {
				log.Printf("recv data failed %v", err)
				return
			}
		}
	}()
	group.Wait()
}

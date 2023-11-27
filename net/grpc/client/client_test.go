package client

import (
	"context"
	"go-lib/net/grpc/pb/grpcdemo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  1.0 * time.Second,
			Multiplier: 1.6,
			Jitter:     0.2,
			MaxDelay:   30 * time.Second,
		},
	}))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}

	for {
		time.Sleep(time.Second * 20)
		cli := grpcdemo.NewGrpcDemoClient(conn)
		result, err := cli.Call(context.Background(), &grpcdemo.NoticeReaderReq{
			Msg:       "",
			NoticeWay: &grpcdemo.NoticeReaderReq_Email{Email: "test"},
		})
		if err != nil {
			log.Printf("Call  failed %v", err)
		}
		log.Printf("================================result================================== %v", result)
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
	result, err := cli.Call(context.Background(), &grpcdemo.NoticeReaderReq{
		Msg:       "",
		NoticeWay: &grpcdemo.NoticeReaderReq_Email{Email: "test"},
	})
	if err != nil {
		log.Printf("Call  failed %v", err)
		return
	}
	log.Printf("result = %v", result)
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
	result, err := cli.Call(context.Background(), &grpcdemo.NoticeReaderReq{
		Msg:       "",
		NoticeWay: &grpcdemo.NoticeReaderReq_Email{Email: "test"},
	})
	if err != nil {
		log.Printf("Call  failed %v", err)
		return
	}
	log.Printf("result = %v", result)
}

package client

import (
	"context"
	"go-lib/net/grpc/pb/grpcdemo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

	time.Sleep(10 * time.Second)
	result, err := cli.UnaryCall(context.Background(), &grpcdemo.NoticeReaderReq{
		Msg:       "",
		NoticeWay: &grpcdemo.NoticeReaderReq_Email{Email: "test"},
	})
	if err != nil {
		log.Printf("Call  failed %v", err)

	}

	log.Printf("result = %v", result)
	for {
		time.Sleep(time.Second * 10)
		result, err := cli.UnaryCall(context.Background(), &grpcdemo.NoticeReaderReq{
			Msg:       "",
			NoticeWay: &grpcdemo.NoticeReaderReq_Email{Email: "test"},
		})
		if err != nil {
			log.Printf("Call  failed %v", err)

		}
		log.Printf("============================result =================================== %v", result)
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
	result, err := cli.UnaryCall(context.Background(), &grpcdemo.NoticeReaderReq{
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
	result, err := cli.UnaryCall(context.Background(), &grpcdemo.NoticeReaderReq{
		Msg:       "",
		NoticeWay: &grpcdemo.NoticeReaderReq_Email{Email: "test"},
	})
	if err != nil {
		log.Printf("Call  failed %v", err)
		return
	}
	log.Printf("result = %v", result)
}

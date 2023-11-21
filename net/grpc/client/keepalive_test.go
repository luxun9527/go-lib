package client

import (
	"context"
	"go-lib/net/grpc/pb/grpcdemo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

func TestKeepAlive(t *testing.T) {
	ctx,cancel :=context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithConnectParams(),grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err!=nil{
		log.Printf("DialContext failed %v",err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)
	result, err := cli.Call(context.Background(), &grpcdemo.NoticeReaderReq{
		Msg:       "",
		NoticeWay: &grpcdemo.NoticeReaderReq_Email{Email: "test"},
	})
	if err!=nil{
		log.Printf("Call  failed %v",err)
		return
	}
	log.Printf("result = %v",result)
}

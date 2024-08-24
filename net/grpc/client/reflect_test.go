package client

import (
	"bytes"
	"context"
	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"strings"
	"testing"

	"google.golang.org/grpc"
)

func TestReflection(t *testing.T) {
	// 1. 建立与 gRPC 服务器的连接
	target := "localhost:8899" // 替换为您的 gRPC 服务器地址
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := grpcreflect.NewClientAuto(context.Background(), conn)
	source := grpcurl.DescriptorSourceFromServer(context.Background(), client)
	reader := strings.NewReader(`{"importedMessage":"test"}`)
	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: true,
		IncludeTextSeparator:  true,
		AllowUnknownFields:    true,
	}
	buffer := bytes.NewBuffer(make([]byte, 0, 100))
	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, source, reader, options)
	h := &grpcurl.DefaultEventHandler{
		Out:            buffer,
		Formatter:      formatter,
		VerbosityLevel: 0,
	}

	rpcPath := "grpcdemo.GrpcDemo/DemoImport"
	if err := grpcurl.InvokeRPC(context.TODO(), source, conn, rpcPath, []string{},
		h, rf.Next); err != nil {
		log.Fatalf("Failed to invoke RPC: %v", err)
	}
	data, _ := io.ReadAll(buffer)
	log.Printf("Response: \n%v", string(data))
}

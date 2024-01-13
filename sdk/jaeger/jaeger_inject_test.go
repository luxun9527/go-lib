package main

import (
	"context"
	"go-lib/net/grpc/pb/grpcdemo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"testing"
	"time"
)

/*
演示如何在网络传输中使用OpenTelemetry。

*/

func TestGrpcInject(t *testing.T) {

	tp, err := tracerProvider("http://192.168.2.159:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	// 设置全局的TracerProvider，方便后面使用
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	// 优雅退出
	defer func(ctx context.Context) {

		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)
	// 新建一个tracer
	tr := otel.Tracer(_globalTrace)
	//开启一个span
	ctx, span := tr.Start(ctx, "func1")
	// trace.SpanFromContext(ctx) 可以通过这个函数，获取span
	startGrpcClient(ctx)
	span.End()
	time.Sleep(time.Second * 3)
}
func startGrpcClient(ctx context.Context) {
	spanCtx := trace.SpanContextFromContext(ctx)
	log.Println(spanCtx.TraceID().String())
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	Inject(ctx, otel.GetTextMapPropagator(), &md)
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)
	result, err := cli.UnaryCall(ctx, &emptypb.Empty{})
	if err != nil {
		log.Printf("Call  failed %v", err)
	} else {
		log.Printf("resp %v", result)
	}

}

type metadataSupplier struct {
	metadata *metadata.MD
}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func (s *metadataSupplier) Set(key, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}

	return out
}

// Inject injects cross-cutting concerns from the ctx into the metadata.
func Inject(ctx context.Context, p propagation.TextMapPropagator, metadata *metadata.MD) {
	p.Inject(ctx, &metadataSupplier{
		metadata: metadata,
	})
}

// Extract extracts the metadata from ctx.
func Extract(ctx context.Context, p propagation.TextMapPropagator, metadata *metadata.MD) (
	baggage.Baggage, trace.SpanContext) {
	ctx = p.Extract(ctx, &metadataSupplier{
		metadata: metadata,
	})

	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}

type GrpcDemoServer struct {
	port int32
	grpcdemo.UnimplementedGrpcDemoServer
}

func (c GrpcDemoServer) UnaryCall(ctx context.Context, req *emptypb.Empty) (*grpcdemo.UnaryCallResp, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	bags, spanCtx := Extract(ctx, otel.GetTextMapPropagator(), &md)
	ctx = baggage.ContextWithBaggage(ctx, bags)
	tr := otel.Tracer(_globalTrace)
	ctx, span := tr.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), "grpc-server")
	defer span.End()
	return &grpcdemo.UnaryCallResp{Username: "zhangsan"}, nil
}

func TestStartGrpcServer(t *testing.T) {
	tp, err := tracerProvider("http://192.168.2.159:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	listener, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}
	s := grpc.NewServer()
	grpcdemo.RegisterGrpcDemoServer(s, new(GrpcDemoServer))
	log.Printf("start server at %v", 8899)
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}

package main

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"log"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const (
	service     = "trace-demo" // 服务名
	environment = "production" // 环境
	id          = 1            // id
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		),
		),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)
	return tp, nil
}

const _globalTrace = "global-trace-id"

func TestJaegerBaseUseage(t *testing.T) {

	tp, err := tracerProvider("http://192.168.2.159:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())
	// 设置全局的TracerProvider，方便后面使用
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 新建一个tracer
	tr := otel.Tracer(_globalTrace)
	//开启一个span
	ctx, span := tr.Start(ctx, "func1")
	// trace.SpanFromContext(ctx) 可以通过这个函数，获取span
	defer span.End()
	func2(ctx)

}

func func2(ctx context.Context) {
	spanCtx := trace.SpanContextFromContext(ctx)
	log.Println(spanCtx.TraceID().String())
	tr := otel.Tracer(_globalTrace)
	// 调用一个方法开启一个span
	_, span := tr.Start(ctx, "fun2")
	span.SetAttributes(attribute.Key("func2_key").String("func2_value"))
	defer span.End()
	time.Sleep(time.Millisecond * 300)
	func3(ctx)
}
func func3(ctx context.Context) {
	spanCtx := trace.SpanContextFromContext(ctx)
	log.Println(spanCtx.TraceID().String())
	tr := otel.Tracer(_globalTrace)
	_, span := tr.Start(ctx, "fun3")
	span.SetAttributes(attribute.Key("func3_key").String("func3_value"))
	defer span.End()
}

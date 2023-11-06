package main

import (
	"context"
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
	// Create the Jaeger exporter
	// 创建 Jaeger exporter
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
		)),
	)
	return tp, nil
}

func TestJaeperSingle1(t *testing.T) {
	tp, err := tracerProvider("http://192.168.2.159:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
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
	// trace 上报
	tr := otel.Tracer("component-main")

	ctx, span := tr.Start(ctx, "root")
	defer span.End()

	//ctx, f1 := tr.Start(ctx, "f1")
	//defer f1.End()
	//ctx, f2 := tr.Start(ctx, "f2")
	//defer f2.End()
	// Context 向下传递
	func1(ctx)

}
func func1(ctx context.Context) {
	// Use the global TracerProvider.
	// 使用 全局 TracerProvider
	tr := otel.Tracer("root")
	_, span := tr.Start(ctx, "fun1")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()
	func2(ctx)

}
func func2(ctx context.Context) {
	// Use the global TracerProvider.
	// 使用 全局 TracerProvider
	tr := otel.Tracer("root")
	_, span := tr.Start(ctx, "fun2")
	span.SetAttributes(attribute.Key("testasetsaetast").String("testasetsaetast"))
	defer span.End()

	// Do fun1...
}

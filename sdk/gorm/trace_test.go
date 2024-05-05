package main

import (
	"context"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"log"
	"testing"
)

const (
	service     = "trace-demo1" // 服务名
	environment = "production"  // 环境
	id          = 1             // id
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
func TestGormTrace(t *testing.T) {
	ctx := context.Background()

	tp, err := tracerProvider("http://192.168.2.200:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Use(otelgorm.NewPlugin(
		otelgorm.WithDBName("test"),
		otelgorm.WithTracerProvider(tp),
	)); err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("gormtracer")

	ctx, span := tracer.Start(ctx, "gormtest")
	defer span.End()
	user := &User{
		Username: "test1",
		Age:      1,
		Fav:      "1",
	}
	//INSERT INTO `user` (`username`,`age`,`fav`,`created_at`,`updated_at`) VALUES ('',0,'',1692947238,1692947238)
	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		fmt.Println("create err", err)
	}

	//otelplay.PrintTraceID(ctx)
	tp.Shutdown(ctx)
}

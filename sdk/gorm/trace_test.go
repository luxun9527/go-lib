package main

import (
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"log"
	"testing"
)
const (
	service     = "trace-demo" // 服务名
	environment = "production" // 环境
	id          = 1      // id
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
	tp, err := tracerProvider("http://192.168.11.185:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	plugin := otelgorm.NewPlugin(
		otelgorm.WithDBName("user"),
		otelgorm.WithAttributes(semconv.DBSystemMySQL, attribute.String("db.addr", "192.168.11.185")),
		otelgorm.WithTracerProvider(tp),
	)
	if err := db.Use(plugin); err != nil {
		log.Fatal(err)
	}
	user := &User{
		Username:  "test",
		Age:       1,
		Fav:       "test",
		CreatedAt: 0,
		UpdatedAt: 0,
	}
	//INSERT INTO `user` (`username`,`age`,`fav`,`created_at`,`updated_at`) VALUES ('',0,'',1692947238,1692947238)
	if err := db.Create(user).Error; err != nil {
		fmt.Println("create err", err)
	}
	fmt.Printf("user = %+v", user)
	//INSERT INTO `user` (`fav`,`created_at`,`updated_at`) VALUES ('',1692947405,1692947405)
	if err := db.Select("fav").Create(user).Error; err != nil {
		fmt.Println("create err", err)
	}
	//INSERT INTO `user` (`username`,`age`,`created_at`,`updated_at`,`id`) VALUES ('',0,1692947813,1692947813,14)
	if err := db.Omit("fav").Create(user).Error; err != nil {
		fmt.Println("create err", err)
	}
	//https://gorm.io/zh_CN/docs/models.html
	//官方的文档很全主要是验证一些模糊的地方
	//1、零值会被插入。
	//2、created_at updated_at会被填充当前时间插入。
	//3、插入后会将主键赋值回来。
}

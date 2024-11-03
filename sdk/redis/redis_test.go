package main

import (
	"context"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"log"
	"sync"

	"testing"

	"time"
)

const (
	service     = "trace-demo" // 服务名
	environment = "production" // 环境
	id          = 1            // id
)

const (
	serviceName    = "redis-Jaeger-Demo"
	jaegerEndpoint = "192.168.2.159:14268/api/traces"
)

var tracer = otel.Tracer("redis-demo")

// newJaegerTraceProvider 创建一个 Jaeger Trace Provider

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

// initTracer 初始化 Tracer
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	tp, err := tracerProvider("http://192.168.2.159:14268/api/traces")
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	return tp, nil
}
func doSomething(ctx context.Context, rdb *redis.ClusterClient) error {
	if err := rdb.Set(ctx, "name", "Q1mi", time.Minute).Err(); err != nil {
		return err
	}
	if err := rdb.Set(ctx, "tag", "OTel", time.Minute).Err(); err != nil {
		return err
	}
	var wg sync.WaitGroup
	for range []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val := rdb.Get(ctx, "tag").Val()
			if val != "OTel" {
				log.Printf("%q != %q", val, "OTel")
			}
		}()
	}
	wg.Wait()

	if err := rdb.Del(ctx, "name").Err(); err != nil {
		return err
	}
	if err := rdb.Del(ctx, "tag").Err(); err != nil {
		return err
	}
	log.Println("done!")
	return nil
}
func TestRedisTrace(t *testing.T) {

	ctx := context.Background()

	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// 创建 Redis 集群客户端
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"192.168.2.159:7001", // 分片节点1的地址
			"192.168.2.159:7002", // 分片节点2的地址
			"192.168.2.159:7003", // 分片节点3的地址
		},
		Password: "bingo", // 如果设置了密码，可以在这里配置
	})

	// 启用 tracing
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	// 启用 metrics
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	ctx, span := tracer.Start(ctx, "doSomething")
	defer span.End()

	if err := doSomething(ctx, rdb); err != nil {
		span.RecordError(err) // 记录error
		span.SetStatus(codes.Error, err.Error())
	}
}

// v9 对应的服务端是7.0 v8对应的是是6.0
func TestRedis(t *testing.T) {
	ctx := context.Background()

	// 创建 Redis 集群客户端
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"192.168.2.159:7001", // 分片节点1的地址
			"192.168.2.159:7002", // 分片节点2的地址
			"192.168.2.159:7003", // 分片节点3的地址
		},
		Password: "bingo", // 如果设置了密码，可以在这里配置
	})

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis cluster: %v", err)
	}
	log.Printf("Connected to Redis cluster: %s", pong)

	// 使用示例：设置和获取一个 key
	err = rdb.Set(ctx, "example_key", "Hello, Redis!", 0).Err()
	if err != nil {
		log.Fatalf("Error setting key: %v", err)
	}

	val, err := rdb.Get(ctx, "example_key").Result()
	if err != nil {
		log.Fatalf("Error getting key: %v", err)
	}

	log.Printf("Value of 'example_key': %s", val)

	// 关闭客户端
	if err := rdb.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
	}
}

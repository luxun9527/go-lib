# 程序可观测-链路追踪

## 核心概念

程序可观测：在多服务，服务多实例的情况下，查看程序的运行情况。主要是三个部分**日志，链路追踪，指标采集。**

本文地址 https://github.com/luxun9527/go-lib/tree/master/sdk/jaeger   如果您觉得这篇文章对您有帮助或启发，请帮我点个star呗，您的star就是我更新的动力。

https://studygolang.com/articles/35630

https://foresightnews.pro/article/detail/25449

### opentelemetry

https://github.com/open-telemetry/docs-cn/tree/main 

https://github.com/open-telemetry/docs-cn/blob/main/OT.md

![img](https://cdn.nlark.com/yuque/0/2023/png/12466223/1699262898081-7b9d493a-834e-4ec7-b98b-50ec8266bf01.png)

opentelemerty 定义 logs(日志) traces(链路追踪)Metrics(服务运行指标采集)的规范。

[https://github.com/open-telemetry/docs-cn/blob/main/OT.md#%E5%A4%A7%E4%B8%80%E7%BB%9F](https://github.com/open-telemetry/docs-cn/blob/main/OT.md#大一统)

## 链路追踪 jaeger



https://www.jaegertracing.io/docs/1.6/getting-started/

https://www.cnblogs.com/whuanle/p/14598049.html

https://www.cnblogs.com/timelesszhuang/p/go-jaeger.html 



https://github.com/yurishkuro/opentracing-tutorial/tree/master/go

https://github.com/jaegertracing/jaeger

https://www.jaegertracing.io/docs/1.53/client-libraries/#deprecating-jaeger-clients

https://opentelemetry.io/docs/demo/services/checkout/#traces

由于这个我使用的不算深入，暂时只分享jaeger go sdk的一些用法。如果您有需求想深入了解可以参考go-zero对于trace相关中间件的实现。

#### 基本概念

链路追踪：在分布式系统中，追踪一次请求或一次操作的调用链，主要是通过context的传播来实现。

jaeper是什么 是opentelmerty 关于 traces链路追踪的实现。

以下是 Jaeger 的一些基本概念：(chatgpt写的)

1. **Span（跨度）**：

- - 在 Jaeger 中，一个 Span 代表一个基本的工作单元。它是在一个服务中处理请求的一段时间。每个请求都会被划分为多个 Span，每个 Span 代表请求的一部分。
  - Span 包含了一些关键信息，比如开始和结束时间、操作名称、标签（key-value 对，用于附加元数据）、日志等。

1. **Trace（跟踪）**：

- - Trace 由一系列的 Span 组成，表示一个完整的请求路径。通过 Trace，你可以了解请求从开始到结束的整个过程，包括请求经过的所有服务和组件。
  - 每个 Span 都有一个唯一的标识符，可以用来关联它们，形成一条 Trace。

1. **Trace Context（跟踪上下文）**：

- - Trace Context 是一组信息，包含在请求的头部中，用于在不同的服务和组件之间传递跟踪信息。它包括 Trace ID（跟踪标识）和 Span ID（跨度标识）等信息。

1. **Sampler（采样器）**：

- - 为了避免过多的数据产生，Jaeger 使用采样器来决定是否要记录一个请求的跟踪信息。采样器可以根据一定的规则，例如采样一定比例的请求，以保持跟踪系统的可用性和性能。

1. **Collector（收集器）**：

- - Jaeger Collector 负责接收来自各个服务的跟踪数据，并将其存储在后端存储系统中。它协调着跟踪数据的收集、存储和检索。

1. **Storage Backend（存储后端）**：

- - 存储后端是实际存储跟踪数据的地方。Jaeger 支持多种存储后端，包括 Elasticsearch、Cassandra、Kafka 等。

1. **Instrumentation（仪表化）**：

- - 为了使用 Jaeger，你需要在应用程序代码中添加一些仪表代码（instrumentation）。这些代码用于创建 Span、记录跟踪信息，并与 Jaeger 进行通信。



![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1704272681032-45dd2ea7-49b6-4772-b84d-b1f807c7e38e.png)



#### docker 安装

使用docker 安装jaeger

可视化ui端口为16686

```shell
docker run -d --name jaeger   -e COLLECTOR_ZIPKIN_HTTP_PORT=9411   -p 5775:5775/udp   -p 6831:6831/udp   -p 6832:6832/udp   -p 5778:5778   -p 16686:16686   -p 14268:14268   -p 9411:9411   jaegertracing/all-in-one:1.6
```

#### go sdk

##### 基本用法

jaeger关于go的的sdk,现在使用的比较多的是go.opentelemetry.io这个包的。



```go
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
    //trace同名
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
```

http://192.168.2.159:16686/ 查看

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1705213473166-979f0d14-3f62-4052-a589-ba7aabf0940c.png)

##### http调用中jaeger

核心就是将 span信息注入到请求头中，然后在服务端解析出来

```go
package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

/*
演示如何在网络传输中使用jaeger。

*/

func TestHttpInject(t *testing.T) {

	tp, err := tracerProvider("http://192.168.2.159:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	// 设置全局的TracerProvider，方便后面使用
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())
	ctx := context.Background()
	// 新建一个tracer
	tr := otel.Tracer(_globalTrace)
	//开启一个span
	ctx, span := tr.Start(ctx, "func1")
	// trace.SpanFromContext(ctx) 可以通过这个函数，获取span
	startHttpClient(ctx)
	span.End()
	time.Sleep(time.Second * 3)
}
func startHttpClient(ctx context.Context) {
	spanCtx := trace.SpanContextFromContext(ctx)
	log.Println(spanCtx.TraceID().String())
	// 创建一个http客户端
	client := &http.Client{}
	propagator := otel.GetTextMapPropagator()

	// 发送GET请求
	req, err := http.NewRequest("POST", "http://localhost:8888/trace", bytes.NewBuffer([]byte(`{"code":"1","value":1"}`)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 注入span信息到请求头中
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	response, err := client.Do(req)
	if err != nil {
		log.Panicf("err %v", err)
	}
	defer response.Body.Close()
	// 读取响应的内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应内容
	fmt.Println("Response:", string(body))
}

func TestGinServer(t *testing.T) {
	engine := gin.New()
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	propagator := otel.GetTextMapPropagator()
	tp, err := tracerProvider("http://192.168.2.159:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	engine.POST("/trace", func(c *gin.Context) {
		// 从请求头中提取span信息
		ctx := propagator.Extract(c, propagation.HeaderCarrier(c.Request.Header))
		spanCtx := trace.SpanContextFromContext(ctx)
		log.Println(spanCtx.TraceID().String())
		tr := tp.Tracer(_globalTrace)
		// 开启一个新的span
		_, span := tr.Start(ctx, "server trace func")
		span.SetAttributes(attribute.Key("server_trace_func_key").
			String("server_trace_func_value"))
		span.End()
		c.JSON(200, gin.H{"traceId": spanCtx.TraceID().String()})
	})
	engine.Run(":8888")
}
```

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1705229400577-c515e4fb-2974-400e-bc2d-4631930944e4.png)

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1705228975357-42ceb8a8-8bcf-467b-837c-05cebed8ddc0.png)

##### grpc调用中使用jaeger

核心就是将 span信息注入到metadata中，然后在服务端解析出来 ，下面的示例截取了go-zero关于grpc trace拦截器的代码。如果想详细学习，也可以去参考 go-zero的相关代码，go-zero还有关于中间件如redis这些加jaeger的代码。

```go
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
演示如何在grcp中使用jaeger。

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

	defer tp.Shutdown(context.Background())

	// 新建一个tracer
	tr := otel.Tracer(_globalTrace)
	//开启一个span
	ctx, span := tr.Start(context.Background(), "func1")
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
```



### 
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
	// 注入span信息注入到请求头中
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

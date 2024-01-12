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
演示如何在网络传输中使用OpenTelemetry。

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
		ctx := propagator.Extract(c, propagation.HeaderCarrier(c.Request.Header))
		spanCtx := trace.SpanContextFromContext(ctx)
		log.Println(spanCtx.TraceID().String())
		tr := tp.Tracer(_globalTrace)
		_, span := tr.Start(ctx, "server trace func")
		span.SetAttributes(attribute.Key("server_trace_func_key").String("server_trace_func_value"))
		span.End()
	})
	engine.Run(":8888")
}

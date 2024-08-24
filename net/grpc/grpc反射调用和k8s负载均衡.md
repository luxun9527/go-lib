# grpc反射和k8s负载均衡

https://www.lixueduan.com/posts/grpc/13-loadbalance-on-k8s/



## 使用反射调用grpc接口

### 客户端

使用grpcurl库

```go
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
```

### 服务端

```go
listener, err := net.Listen("tcp", "0.0.0.0:8899")
if err != nil {
    log.Println("net listen err ", err)
    return
}

s := grpc.NewServer()
//开启反射
reflection.Register(s)
grpcdemo.RegisterGrpcDemoServer(s, new(GrpcDemoServer))
log.Printf("start server at %v", 8899)
if err := s.Serve(listener); err != nil {
    log.Println("failed to serve...", err)
    return
}
```

## k8s下grpc的负载均衡

grpc是长连接，k8s下如何实现负载均衡，具体和使用etcd差不多,需要自定义resolver，在每次请求的时候获取

svc下pod的变化，建立连接

关于服务的解析可以参考

https://www.liwenzhou.com/posts/Go/name-resolving-and-load-balancing-in-grpc/

https://www.lixueduan.com/posts/grpc/13-loadbalance-on-k8s/

```go
package main

import (
    "context"
    "flag"
    "github.com/sercand/kuberesolver/v5"
    "go-lib/net/grpc/pb/grpcdemo"
    "go-lib/net/grpc/server"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/reflection"
    "google.golang.org/protobuf/types/known/emptypb"
    "log"
    "net"
    "time"
)

var (
    mode = flag.String("mode", "server", "server or client")
)

func main() {
    flag.Parse()
    log.Println("mode:", *mode)
    if *mode == "client" {
        initClient()
    } else {
        initServer()
    }

}
func initServer() {
    listener, err := net.Listen("tcp", "0.0.0.0:8899")
    if err != nil {
        log.Println("net listen err ", err)
        return
    }

    s := grpc.NewServer()
    //开启反射
    reflection.Register(s)
    grpcdemo.RegisterGrpcDemoServer(s, new(server.GrpcDemoServer))
    log.Printf("start server at %v", 8899)
    if err := s.Serve(listener); err != nil {
        log.Println("failed to serve...", err)
        return
    }
}
func initClient() {
    kuberesolver.RegisterInCluster()

    conn, err := grpc.Dial("kubernetes:///grpcdemoserver-service.default:8899",
                           grpc.WithTransportCredentials(insecure.NewCredentials()),
                           grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
                          )
    if err != nil {
        log.Printf("DialContext failed %v", err)
        return
    }
    cli := grpcdemo.NewGrpcDemoClient(conn)

    for {
        time.Sleep(time.Second * 5)
        result, err := cli.UnaryCall(context.Background(), &emptypb.Empty{})
        if err != nil {
            log.Printf("Call  failed %v", err)
        } else {
            log.Printf("resp %v", result)
        }
    }
}
```
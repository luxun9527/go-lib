# grpc

refer

https://www.liwenzhou.com/posts/Go/name-resolving-and-load-balancing-in-grpc/

https://github.com/win5do/go-microservice-demo/blob/main/docs/sections/grpc-lb.md

https://www.cnblogs.com/FireworksEasyCool/p/12912839.html

https://cloud.tencent.com/developer/article/2136435

https://www.liwenzhou.com/posts/Go/name-resolving-and-load-balancing-in-grpc/

介绍grpc常用用法。

本文代码地址https://github.com/luxun9527/go-lib/tree/master/net/grpc ，您的star就是我更新的动力。

## 基本概念

1、rpc是什么

RPC（Remote Procedure Call）远程过程调用协议，可以理解为一种抽象的协议（类似接口），定义了程序远程调用另一个程序要实现那些东西。

总体来说有以下几步

- 规定传输协议
- 建立连接。
- 发送方序列化传输的数据，接收方反序列化收到的数据。

 restfull(http+json) ，jsonrpc(tcp+json) grpc(http2+protobuf) 这些都可以叫做rpc。

2、grpc是什么

grpc一种rpc协议的实现。grpc使用http2作为传输协议，protobuf作为序列化协议，go比较有名的还有rpcx。

## protobuf

https://github.com/protocolbuffers/protobuf/releases protoc 下载

```go
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0 
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.0
go install  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

下面是protobuf文件包含了protobuf的常用用法，包括message，service定义，如何导入其他的message ,-I flag，protobuf相关命令的用法。

其他详细的message用法可以参考https://alovn.cn/docs/protobuf/proto3/json.html

google官方定义的message https://github.com/protocolbuffers/protobuf/releases下载protoc的时候里面有，具体如下。

https://github.com/luxun9527/go-lib/tree/master/net/grpc/pb/googleapis/google/protobuf 

https://github.com/luxun9527/go-lib/tree/master/net/grpc/pb/googleapis/google/api

```go
syntax = "proto3";
//当别人导入这个protobuf文件，使用的包名 如 google/protobuf/empty.proto 定义的就是 package google.protobuf,我们要使用这个文件中message 使用方法为 package.Message
//如google.protobuf(包名).Empty(message)
package grpcdemo;

//go_package = "./grpcdemo;grpcdemo"; ./grpcdemo表示生成的文件的位置和生成命令指定的生成位置,一起决定最后生成文件的位置 grpcdemo表示生成pb文件的包名
option go_package = "./grpcdemo;grpcdemo";


//导入其他protobuf 导入我们自定义的protobuf 需要和protoc 命令 -I参数组成完整的导入路径。例如，导入google/protobuf/empty.proto需要指定 -I./pb/googleapis
import "google/protobuf/empty.proto";

//导入我们自定义的protobuf 需要和  protoc -I参数组成完整的导入路径。
import "grpcdemo/folder/imported.proto";

//特殊情况当被导入的proto和我们是同一级的时候。可以不使用package.Message的形式 直接使用即可，CustomMessage
import "grpcdemo/custom.proto";

import "google/api/annotations.proto";

service GrpcDemo {
    //grpc 4种调用类型

    //Unary RPC （一元RPC）
    rpc UnaryCall(NoticeReaderReq)returns(google.protobuf.Empty);

    //Unary RPC （一元RPC）
    rpc DemoImport(grpcdemo.folder.ImportedMessage)returns(CustomMessage);

    //Client Streaming RPC （ 客户端流式RPC）
    rpc PushData(stream Empty) returns(Data);

    //Server Streaming RPC （ 服务器流式RPC）
    rpc FetchData(Empty) returns(stream Data);

    //Bidirectional Streaming RPC （双向流式RPC）
    rpc Exchange(stream Req) returns(stream Resp);
    //grpc-gateway调用
    rpc CallGrpcGateway(NoticeReaderReq)returns(NoticeReaderResp){
        option (google.api.http) = {
            post: "/v1/call"
            body:"*"
        };
    }
}
message Req{
    string firstName =1;
    optional string age=2;
}
message Resp{
    string lastName=1;
    Gender gender=2;
}
//枚举类
enum Gender{
    Unknown =0;
    Male=1;
    Female=2;
}

message Empty{}

message Data{
    string uid =1;
    string topic=2;
    bytes data=3;
}
message NoticeReaderResp{
    string fav_book=4;//最爱的书
}

// protobuf oneof的用法。
message NoticeReaderReq{
    string msg = 1;

    oneof notice_way{
        string email = 2;
        string phone = 3;
    }
}

/*
 一个pb文件可以定义多个service
*/
service GrpcGatewayDemo {
    rpc CallGrpcGatewayDemo(NoticeReaderReq)returns(NoticeReaderResp){
        option (google.api.http) = {
            post: "/v1/gateway"
            body:"*"
        };
    }
}
```

```makefile
.PHONY: proto
proto:
    protoc   -I./pb -I./pb/googleapis  --grpc-gateway_out=./pb \
          --openapiv2_out=. \
          --openapiv2_opt json_names_for_fields=false\
          --openapiv2_opt generate_unbound_methods=true \
          --openapiv2_opt output_format=yaml \
          --grpc-gateway_opt generate_unbound_methods=true \
          --grpc-gateway_opt logtostderr=true \
            --go_out=./pb --go-grpc_out=./pb grpcdemo/grpcdemo.proto
    protoc   -I./pb --go_out=./pb --openapiv2_out=.  grpcdemo/custom.proto
    protoc   -I./pb --go_out=../../../ --openapiv2_out=. grpcdemo/folder/imported.proto
	
	#--openapiv2_opt 定义参数 https://github.com/grpc-ecosystem/grpc-gateway/blob/main/protoc-gen-openapiv2/main.go
   #--grpc-gateway_opt定义参数 https://github.com/grpc-ecosystem/grpc-gateway/blob/main/protoc-gen-grpc-gateway/main.go
	#-I 表示要从哪里开始找protobuf文件。--openapiv2_out生成的swagger文件放在哪
	# --grpc-gateway_out=./pb 生成的grpcgateway文件放在哪 --go_out=./pb  --go-grpc_out=./pb 也是一样的。
	
    # 如果有指定-I的操作 要么一定要有一个-I参数能够和proto文件构成完整导入路径，例如 -I./pb 和 grpcdemo/custom.proto 构成了custom.proto文件的完整路径

```



## grpc四种模式

1、四种模式根据前面protobuf中的定义，生成的代码，基本的使用示例如下。

服务端

```go
package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
)

type GrpcDemoServer struct {
	port int32
	grpcdemo.UnimplementedGrpcDemoServer
}

func (c GrpcDemoServer) UnaryCall(ctx context.Context, req *emptypb.Empty) (*grpcdemo.UnaryCallResp, error) {
	log.Printf("port is %v", c.port)
	return &grpcdemo.UnaryCallResp{Username: "zhangsan"}, nil
}
func (c GrpcDemoServer) DemoImport(ctx context.Context, req *folder.ImportedMessage) (*grpcdemo.CustomMessage, error) {
	log.Printf("port is %v", c.port)
	return &grpcdemo.CustomMessage{
		CustomMessage: fmt.Sprintf("%v", c.port),
	}, nil
}
func (GrpcDemoServer) PushData(c grpcdemo.GrpcDemo_PushDataServer) error {
	for {
		data, err := c.Recv()
		if err != nil {
			log.Printf("err %v", err)
			return err
		}
		log.Printf("recv data %v", data)
	}

}
func (GrpcDemoServer) FetchData(req *grpcdemo.FetchDataReq, c grpcdemo.GrpcDemo_FetchDataServer) error {
	for i := 0; i < 10; i++ {
		if err := c.Send(&grpcdemo.FetchDataResp{
			FavBook: "book",
		}); err != nil {
			log.Printf("err %v", err)
			return err
		}
	}
	return nil
}
func (GrpcDemoServer) Exchange(c grpcdemo.GrpcDemo_ExchangeServer) error {
	g := sync.WaitGroup{}

	g.Add(2)
	go func() {
		defer g.Done()
		for {
			data, err := c.Recv()
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("exchange recv message %v", data)
		}

	}()
	go func() {
		defer g.Done()
		for {
			if err := c.Send(&grpcdemo.ExchangeResp{LastName: "test"}); err != nil {
				if err != nil {
					log.Println(err)
					return
				}
			}
		}

	}()
	g.Wait()
	return nil
}

func (GrpcDemoServer) CallGrpcGateway(ctx context.Context, req *grpcdemo.CallGrpcGatewayReq) (*grpcdemo.CallGrpcGatewayResp, error) {
	log.Printf("recv message %v", req.Config)
	name := req.Config["name"]
	switch name {
	case "zhangsan":
		return nil, status.Error(codes.NotFound, "not found")
	case "lisi":
		return nil, errors.New("this is custom error")

	}
	return &grpcdemo.CallGrpcGatewayResp{Config: req.Config}, nil
}

type GrpcGatewayDemo struct {
	grpcdemo.GrpcGatewayDemoServer
}

func (GrpcGatewayDemo) CallGrpcGatewayDemo(ctx context.Context, req *grpcdemo.CallGrpcGatewayDemoReq) (*grpcdemo.CallGrpcGatewayDemoResp, error) {

	return &grpcdemo.CallGrpcGatewayDemoResp{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

func TestServer(t *testing.T) {
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

func TestGrpcGateWayServer(t *testing.T) {
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:8899")
		if err != nil {
			log.Println("net listen err ", err)
			return
		}
		s := grpc.NewServer()
		grpcdemo.RegisterGrpcDemoServer(s, new(GrpcDemoServer))
		grpcdemo.RegisterGrpcGatewayDemoServer(s, new(GrpcGatewayDemo))
		if err := s.Serve(listener); err != nil {
			log.Println("failed to serve...", err)
			return
		}
	}()
	conn, err := grpc.Dial(
		"127.0.0.1:8899",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Panic("dail proxy grpc serve failed ", zap.Error(err))
	}

	gwmux := runtime.NewServeMux()

	if err = grpcdemo.RegisterGrpcDemoHandler(context.Background(), gwmux, conn); err != nil {
		log.Panicf("Failed to register gateway %v", err)
	}
	if err = grpcdemo.RegisterGrpcGatewayDemoHandler(context.Background(), gwmux, conn); err != nil {
		log.Panicf("Failed to register gateway %v", err)
	}

	gwServer := &http.Server{
		Addr:    ":10080",
		Handler: gwmux,
	}
	if err := gwServer.ListenAndServe(); err != nil {
		log.Panic("init proxy http serve failed err", zap.Error(err))
	}

}


```

客户端

```go
package client

import (
	"context"
	"go-lib/net/grpc/pb/grpcdemo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sync"
	"testing"
	"time"
)

func TestClient(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)

	for {
		time.Sleep(time.Second * 10)
		result, err := cli.UnaryCall(context.Background(), &emptypb.Empty{})
		if err != nil {
			log.Printf("Call  failed %v", err)
		} else {
			log.Printf("resp %v", result)
		}
	}
}

func TestPush(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)

	c, err := cli.PushData(context.Background())
	if err != nil {
		log.Printf("get pushdata connection failed %v", err)
		return
	}
	for i := 0; i < 10; i++ {
		if err := c.Send(&grpcdemo.PushDataReq{Foo: "foo"}); err != nil {
			log.Printf("push data failed %v", err)
			return
		}
	}

}

func TestFetchData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)
	c, err := cli.FetchData(context.Background(), &grpcdemo.FetchDataReq{
		Msg:       "",
		NoticeWay: &grpcdemo.FetchDataReq_Email{Email: "test"},
	})
	if err != nil {
		log.Printf("get fetchdata connection failed %v", err)
		return
	}

	for {
		data, err := c.Recv()
		if err != nil {
			log.Printf("recv data failed %v", err)
			return
		}
		log.Printf("data =%v", data)
	}

}
func TestExchangeData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:8899", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("DialContext failed %v", err)
		return
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)
	c, err := cli.Exchange(context.Background())
	if err != nil {
		log.Printf("get Exchangedata connection failed %v", err)
		return
	}
	group := sync.WaitGroup{}
	group.Add(2)
	go func() {
		for {
			defer group.Done()
			data, err := c.Recv()
			if err != nil {
				log.Printf("recv data failed %v", err)
				return
			}
			log.Printf("data =%v", data)
		}
	}()
	go func() {
		for {
			defer group.Done()
			var age = "12"
			err := c.Send(&grpcdemo.ExchangeReq{
				FirstName: "zhangsan",
				Age:       &age,
			})
			if err != nil {
				log.Printf("recv data failed %v", err)
				return
			}
		}
	}()
	group.Wait()
}

```

## grpc自定义target解析

场景 ：我们有多个服务端实例怎么去实现负载均衡，dial直接使用地址只能连接一个，这时候可以使用自定义target解析，并使用round_robin负载均衡策略来实现。

**客户端**

```go
package client

import (
    "context"
    "go-lib/net/grpc/pb/grpcdemo"
    "go-lib/net/grpc/pb/grpcdemo/folder"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/resolver"
    "log"
    "testing"
    "time"
)

/*
refer
https://www.liwenzhou.com/posts/Go/name-resolving-and-load-balancing-in-grpc/
https://github.com/win5do/go-microservice-demo/blob/main/docs/sections/grpc-lb.md
https://www.cnblogs.com/FireworksEasyCool/p/12912839.html
*/
// 自定义name resolver

var (
    _customScheme   = "dns"
    _customEndpoint = "xxx.xxx.com"
    _addrs          = []string{"127.0.0.1:8898", "127.0.0.1:8899"}
)

// customResolver 自定义name resolver，实现Resolver接口
type customResolver struct {
}

func (r *customResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*customResolver) Close() {}

// customBuilder 需实现 Builder 接口
type customBuilder struct {
    target     resolver.Target
    cc         resolver.ClientConn
    addrsStore map[string][]string
}

func (builder *customBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

    builder.target = target
    builder.cc = cc
    builder.addrsStore = map[string][]string{_customEndpoint: _addrs}
    addresses := builder.addrsStore[target.Endpoint()]
    r := make([]resolver.Address, 0, 2)
    for _, v := range addresses {
       a := resolver.Address{
          Addr:       v,
          ServerName: "",
       }
       r = append(r, a)
    }
    //最核心的就是执行这个函数更新地址列表。
    if err := builder.cc.UpdateState(resolver.State{Addresses: r}); err != nil {
       return nil, err
    }

    go func() {
        //20秒后更新连接新增一个地址。
       time.Sleep(time.Second * 20)
       builder.updateConn()
    }()
    return &customResolver{}, nil
}

// 执行UpdateState更新连接
func (builder *customBuilder) updateConn() {
    addresses := []resolver.Address{{
       Addr:               "127.0.0.1:8898",
    }, {
       Addr:               "127.0.0.1:8899",
    }, {
       Addr:               "127.0.0.1:8897",
    }}
    if err := builder.cc.UpdateState(resolver.State{Addresses: addresses}); err != nil {
       log.Printf("update state failed %v", err)
    }
}
func (*customBuilder) Scheme() string { return _customScheme }

func TestResolverClientTest(t *testing.T) {
    resolverBuilder := &customBuilder{}
    conn, err := grpc.Dial(
       "dns:///xxx.xxx.com",
       grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
       grpc.WithTransportCredentials(insecure.NewCredentials()),
       grpc.WithResolvers(resolverBuilder),
    )
   // 	resolver.Register(&customBuilder{}) 也可以使用全局注册的方式

    /*
    pick_first
尝试连接到第一个地址，如果连接成功，则将其用于所有RPC，如果连接失败，则尝试下一个地址（并继续这样做，直到一个连接成功）。
round_robin
连接到所有地址，并依次向每个后端发送一个RPC。例如，第一个RPC将发送到backend-1，第二个RPC将发送到backend-2，第三个RPC将再次发送到backend-1。
    */
    if err != nil {
       log.Printf("dial connection failed err =%v", err)
    }

    cli := grpcdemo.NewGrpcDemoClient(conn)
    for i := 0; i < 200; i++ {
       resp, err := cli.DemoImport(context.Background(), &folder.ImportedMessage{
          ImportedMessage: "test",
       })
       if err != nil {
          log.Printf("call failed err %v", err)
       } else {
          log.Printf("resp %v", resp)
       }
       time.Sleep(time.Second)
    }

}
2023/12/04 22:48:58 resp custom_message:"8899"
2023/12/04 22:48:59 resp custom_message:"8898"
2023/12/04 22:49:00 resp custom_message:"8899"
2023/12/04 22:49:01 resp custom_message:"8898"
2023/12/04 22:49:02 resp custom_message:"8899"
2023/12/04 22:49:03 resp custom_message:"8898"
2023/12/04 22:49:04 resp custom_message:"8899"
2023/12/04 22:49:05 resp custom_message:"8898"
2023/12/04 22:49:06 resp custom_message:"8899"
2023/12/04 22:49:07 resp custom_message:"8898"
2023/12/04 22:49:08 resp custom_message:"8899"
2023/12/04 22:49:09 resp custom_message:"8898"
2023/12/04 22:49:10 resp custom_message:"8899"
2023/12/04 22:49:11 resp custom_message:"8898"
2023/12/04 22:49:12 resp custom_message:"8899"
2023/12/04 22:49:13 resp custom_message:"8898"
2023/12/04 22:49:14 resp custom_message:"8899"
2023/12/04 22:49:15 resp custom_message:"8899"
2023/12/04 22:49:16 resp custom_message:"8897"
2023/12/04 22:49:17 resp custom_message:"8898"
2023/12/04 22:49:18 resp custom_message:"8899"
2023/12/04 22:49:19 resp custom_message:"8897"
2023/12/04 22:49:20 resp custom_message:"8898"
2023/12/04 22:49:21 resp custom_message:"8899"
2023/12/04 22:49:22 resp custom_message:"8897"
2023/12/04 22:49:23 resp custom_message:"8898"
2023/12/04 22:49:24 resp custom_message:"8899"
```

**服务端**

```go
package server

import (
    "go-lib/net/grpc/pb/grpcdemo"
    "google.golang.org/grpc"
    "log"
    "net"
    "sync"
    "testing"
)

func TestResolverServer(t *testing.T) {
    w := sync.WaitGroup{}
    w.Add(3)
    go func() {
       defer w.Done()
       listener, err := net.Listen("tcp", "0.0.0.0:8898")
       if err != nil {
          log.Println("net listen err ", err)
          return
       }
       s := grpc.NewServer()
       grpcdemo.RegisterGrpcDemoServer(s, &GrpcDemoServer{
          port:                        8898,
          UnimplementedGrpcDemoServer: grpcdemo.UnimplementedGrpcDemoServer{},
       })
       log.Printf("start server at %v", 8898)
       if err := s.Serve(listener); err != nil {
          log.Println("failed to serve...", err)
          return
       }
    }()

    go func() {
       defer w.Done()
       listener, err := net.Listen("tcp", "0.0.0.0:8899")
       if err != nil {
          log.Println("net listen err ", err)
          return
       }
       s := grpc.NewServer()
       grpcdemo.RegisterGrpcDemoServer(s, &GrpcDemoServer{
          port:                        8899,
          UnimplementedGrpcDemoServer: grpcdemo.UnimplementedGrpcDemoServer{},
       })
       log.Printf("start server at %v", 8899)
       if err := s.Serve(listener); err != nil {
          log.Println("failed to serve...", err)
          return
       }

    }()
    go func() {
       defer w.Done()
       listener, err := net.Listen("tcp", "0.0.0.0:8897")
       if err != nil {
          log.Println("net listen err ", err)
          return
       }
       s := grpc.NewServer()
       grpcdemo.RegisterGrpcDemoServer(s, &GrpcDemoServer{
          port:                        8897,
          UnimplementedGrpcDemoServer: grpcdemo.UnimplementedGrpcDemoServer{},
       })
       log.Printf("start server at %v", 8897)
       if err := s.Serve(listener); err != nil {
          log.Println("failed to serve...", err)
          return
       }

    }()
    w.Wait()
}
```

## grpc自定义负载均衡

grpc 的客户端的负载均衡可以通过自定义客户端balance来实现。

自定义负载均衡日常使用基本用不到，如果有兴趣可以参考go-zero的实现。

https://www.cnblogs.com/kevinwan/p/16571213.html

https://github.com/zeromicro/go-zero/blob/master/zrpc/internal/balancer/p2c/p2c.go

grpc也可以服务端来实现负载均衡具体可以参考。·

https://github.com/win5do/go-microservice-demo/blob/main/docs/sections/grpc-lb.md
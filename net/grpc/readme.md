refer https://github.com/mofei1/grpc-demo

# rpc

定义rpc 远程调用的一种实现方案，他定义了一个程序远程调用另一个程序的实现方案

包括 1、接口的定义 2、数据的序列化和反序列化 3、数据的传输

# protobuf

grpc的数据序列化和反序列的协议

安装对应的go和grpc 插件

```shell
$ go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.7.0\
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.7.0 \
    google.golang.org/protobuf/cmd/protoc-gen-go@v1.26.0 \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
```

protobuf 的安装直接在github上下就可以。

 protoc 版本libprotoc 3.19.4

protoc-gen-go v1.26.0

protoc-gen-go-grpc 1.1.0

```protobuf
syntax = "proto3";
package grpc;
//./hello 文件生成的路径 user表示包名
option go_package = "./hello;hellopb";

message HelloRequest {
    string name=1;
    int32 age =2;
}

message HelloResponse {
    string reply = 1;
    string second_name=2;
}

service HelloService {
    //方法
    rpc processOrder(HelloRequest) returns(HelloResponse);
}
```



使用的命令

```shell
	protoc  -Ipb -I. --go_out=./ --go-grpc_out=./ pb/*.proto
```

options 介绍

-I   --proto_path 表示 proto文件引入别的protoc文件查找的路径 或 指定protoc要使用的proto文件的路径

--go_out  go文件生成的路径和文件中的`option go_package = "./hello;hellopb";`构成了go文件的生成路径 上面这个命令的生成路径为`./pbpath/hello/hell.pb.go` ./pbpath 必须要存在，hello 这个路径会自动生成。 

--go-grpc_out=./  和--go_out 同样的用法生成grpc代码 文件名为`hello_grpc.pb.go`

最佳实践， --go_out=./ --go-grpc_out=./ 使用当前路径，通过文件中的option 指定文件的路径。
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
)

type GrpcDemoServer struct {
	grpcdemo.UnimplementedGrpcDemoServer
}

func (GrpcDemoServer) Call(ctx context.Context, req *grpcdemo.NoticeReaderReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Call not implemented")
}
func (GrpcDemoServer) DemoImport(ctx context.Context, req *folder.ImportedMessage) (*grpcdemo.CustomMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DemoImport not implemented")
}
func (GrpcDemoServer) PushData(c grpcdemo.GrpcDemo_PushDataServer) error {
	for {
		data, err := c.Recv()
		if err != nil {
			log.Printf("err %v", err)
			return err
		}
		log.Println(data)
	}

}
func (GrpcDemoServer) FetchData(req *grpcdemo.Empty, c grpcdemo.GrpcDemo_FetchDataServer) error {
	for i := 0; i < 10; i++ {
		if err := c.Send(&grpcdemo.Data{}); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
func (GrpcDemoServer) Exchange(c grpcdemo.GrpcDemo_ExchangeServer) error {
	g := sync.WaitGroup{}

	g.Add(2)
	go func() {
		for {
			req, err := c.Recv()
			if err != nil {
				log.Println(err)
			}
			log.Println(req)
		}
		g.Done()

	}()
	go func() {
		for {
			if err := c.Send(&grpcdemo.Resp{LastName: "test"}); err != nil {
				if err != nil {
					log.Println(err)
				}
			}
		}
		g.Done()
	}()
	g.Wait()
	return nil
}

func (GrpcDemoServer) CallGrpcGateway(ctx context.Context, req *grpcdemo.NoticeReaderReq) (*grpcdemo.NoticeReaderResp, error) {
	log.Println(req)
	switch req.Msg {
	case "1":
		return nil,status.Error(codes.NotFound,"not found")
	case "2":
		return nil,fmt.Errorf("custom error")

	}
	return &grpcdemo.NoticeReaderResp{FavBook: ""}, nil
}


type GrpcGatewayDemo struct {
	grpcdemo.GrpcGatewayDemoServer
}

func (GrpcGatewayDemo)CallGrpcGatewayDemo(ctx context.Context,req *grpcdemo.NoticeReaderReq) (*grpcdemo.NoticeReaderResp, error){
	switch req.Msg {
	case "1":
		return nil,status.Error(codes.NotFound,"CallGrpcGatewayDemo not found")
	case "2":
		return nil,fmt.Errorf("CallGrpcGatewayDemo custom error")

	}
	return &grpcdemo.NoticeReaderResp{FavBook: "test11"}, nil
}

type Response struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func TestMarshal(t *testing.T) {
	r := Response{
		Code: 0,
		Msg:  "success",
		Data: struct {
		}{},
	}
	d, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}
	log.Println([]byte{':'})
	indexByte := bytes.LastIndexByte(d, ':')
	log.Println(indexByte)

	for _, v := range d {
		fmt.Printf("%v,", v)
	}

}
func TestServer(t *testing.T) {
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
}

var (
	_responsePrefix = []byte{123,34,99,111,100,101,34,58,48,44,34,109,115,103,34,58,34,115,117,99,99,101,115,115,34,44,34,100,97,116,97,34,58}
	_responseSuffix = []byte{125}
	_responsePrefixLen = len(_responsePrefix)
	_responseSuffixLen = len(_responseSuffix)
	_empty = struct {}{}
)


func (j *JSONPbWrap) Marshal(v interface{}) ([]byte, error) {
	data, err := j.JSONPb.Marshal(v)
	if err != nil {
		return nil, err
	}
	//给返回的数据加上加上我们自定义的一些数据，使其统一格式。
	extra := _responsePrefixLen + _responseSuffixLen
	size :=len(data)+extra
	d := make([]byte, len(data)+extra)
	copy(d, _responsePrefix)
	copy(d[_responsePrefixLen:], data[:])
	copy(d[size-1:], _responseSuffix)
	return d, nil
}

type JSONPbWrap struct {
	runtime.JSONPb
}

func TestGrpcGateWay(t *testing.T) {
	conn, err := grpc.Dial(
		"127.0.0.1:8899",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Panic("dail proxy grpc serve failed ", zap.Error(err))
	}
	marshaler := &runtime.HTTPBodyMarshaler{
		Marshaler: &JSONPbWrap{runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
				UseProtoNames:   true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
		},
	}
	m := runtime.WithMarshalerOption("application/json", marshaler)
	errorHandler := runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
		// return Internal when Marshal failed
		var r Response
		s, ok := status.FromError(err)
		if !ok {
			r = Response{
				Code: int32(codes.Unknown),
				Msg:  err.Error(),
				Data: _empty,
			}
		}
		r = Response{
			Code: int32(s.Code()),
			Msg:  s.Message(),
			Data: _empty,
		}
		buf, _ := json.Marshal(r)
		writer.Write(buf)
	})
	gwmux := runtime.NewServeMux(errorHandler,m)

	if err = grpcdemo.RegisterGrpcDemoHandler(context.Background(), gwmux, conn); err != nil {
		log.Panic("Failed to register gateway ", zap.Error(err))
	}
	if err = grpcdemo.RegisterGrpcGatewayDemoHandler(context.Background(), gwmux, conn); err != nil {
		log.Panic("Failed to register gateway ", zap.Error(err))
	}

	gwServer := &http.Server{
		Addr:    ":10080",
		Handler: gwmux,
	}

	go func() {
		if err := gwServer.ListenAndServe(); err != nil {
			log.Panic("init proxy http serve failed err", zap.Error(err))
		}
	}()
	select {}
}

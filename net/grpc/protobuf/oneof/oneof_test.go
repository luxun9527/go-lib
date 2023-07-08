package oneof

import (
	"context"
	"encoding/json"
	oneof "go-lib/net/grpc/protobuf/oneof/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := grpc.Dial("0.0.0.0:9090", grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect.", err)
		return
	}
	defer conn.Close()

	helloClient := oneof.NewHelloServiceClient(conn)
	helloClient.SayHello(context.Background(), &oneof.NoticeReaderRequest{
		Msg:       "",
		NoticeWay: &oneof.NoticeReaderRequest_Email{Email: ""},
	})

}
func TestMashal(t *testing.T) {
	msg := &oneof.NoticeReaderRequest{
		Msg:       "",
		NoticeWay: &oneof.NoticeReaderRequest_Email{Email: "xxxx@qq.com"},
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	var m oneof.NoticeReaderRequest
	if err := proto.Unmarshal(data, &m); err != nil {
		log.Println(err)
	}
	log.Println(m.NoticeWay)
	person := &Person{Name: "122"}

	data, err = json.Marshal(person)
	var p interface{}
	err = json.Unmarshal(data, &p)
	if err != nil {
		log.Println(err)
	}
	log.Println(p)
}

type Person struct {
	Name string
}

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	oneof.RegisterHelloServiceServer(s, new(Hello))
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}

type Hello struct {
	oneof.UnimplementedHelloServiceServer
}

func (h Hello) SayHello(ctx context.Context, req *oneof.NoticeReaderRequest) (*oneof.Empty, error) {
	switch v := req.NoticeWay.(type) {
	case *oneof.NoticeReaderRequest_Email:
		log.Println(v)
	case *oneof.NoticeReaderRequest_Phone:
		//noticeWithPhone(v)
	}
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}

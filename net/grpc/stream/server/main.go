package main

import (
	"fmt"
	hellopb "go-lib/net/grpc/stream/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type Hello struct {
	hellopb.UnimplementedHelloStreamServer
}

func (Hello) FetchData(req *hellopb.Empty, data hellopb.HelloStream_FetchDataServer) error {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		data.Send(&hellopb.Data{
			Uid:   "1",
			Topic: "1",
			Data:  []byte("abcdefj"),
		})

		data.SendMsg(&hellopb.Data{
			Uid:   "1",
			Topic: "1",
			Data:  []byte("abcdefj"),
		})
	}
	return nil
}
func (Hello) Exchange(data hellopb.HelloStream_ExchangeServer) error {
	go func() {
		for {
			d, err := data.Recv()
			if err == io.EOF {
				fmt.Printf("eof\n")
				break
			}
			fmt.Printf("data = %+v\n", d)
			data.Send(&hellopb.Resp{LastName: "san"})
		}
	}()

	return nil

}
func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		log.Println("net listen err ", err)
		return
	}

	s := grpc.NewServer()
	hellopb.RegisterHelloStreamServer(s, new(Hello))
	if err := s.Serve(listener); err != nil {
		log.Println("failed to serve...", err)
		return
	}
}

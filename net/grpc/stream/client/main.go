package main

import (
	"context"
	"fmt"
	hellopb "go-lib/net/grpc/stream/pb"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	conn, err := grpc.DialContext(timeout,"0.0.0.0:9090", grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect.", err)
		return
	}
	defer conn.Close()

	helloClient := hellopb.NewHelloStreamClient(conn)
	c := Client{helloClient}
	c.fetchData()
	//c.exchange()
}

type Client struct {
	hellopb.HelloStreamClient
}

func (c *Client) fetchData() {
	ctx := context.Background()
	resp, err := c.FetchData(ctx, &hellopb.Empty{})
	if err != nil {
		fmt.Printf("error = %v", err)
		return
	}
	for {
		msg, err := resp.Recv()

		if err != nil {
			log.Println("err", err)
			break
		}

		fmt.Printf("message = %+v\n", msg)
	}
}
func (c *Client) exchange() {
	ctx := context.Background()
	resp, err := c.Exchange(ctx)
	if err != nil {
		fmt.Printf("error = %v", err)
		return
	}
	for {

		d, err := resp.Recv()
		if err != nil {
			fmt.Printf("err = %v\n", err)
			break
		}
		fmt.Printf("data = %+v\n", d)
		if err := resp.Send(&hellopb.Req{FirstName: "zhang"}); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}

}

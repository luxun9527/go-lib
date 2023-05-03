package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*20)
	defer cancel()
	err, data := doRequestWithContext("http://www.baidu.com/", ctx)
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}
	fmt.Printf("data = %v\n", string(data))

}
func doRequestWithContext(url string, ctx context.Context) (error, []byte) {

	dc, ec := make(chan []byte), make(chan error)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ec <- errors.New(err.(string))
			}
		}()
		resp, err := http.Get(url)
		if err != nil {
			ec <- err
			return
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ec <- err
			return
		}
		dc <- data

	}()
	select {
	case <-ctx.Done():
		return ctx.Err(), nil
	case data := <-dc:
		return nil, data
	case err := <-ec:
		return err, nil
	}

}

package main

import (
	"bytes"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kr/pretty"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//https://github.com/go-resty/resty
func main() {

}

type RWC struct {
	data []byte
	file *os.File
}

//func(r *RWC) Read(p []byte) (n int, err error)  {
//	read, err := r.file.Read(p)
//
//	return 0, err
//}
//
//func(r *RWC) Write() (n int, err error)  {
//
//	buf :=make([]byte,1024*1024*4)
//	data, err := r.Read(buf)
//
//	return 0, err
//}
//func(r *RWC) Close() error   {
//	return nil
//}
//实时读写
func client() {

	resp, err := http.Post("http://127.0.0.1:9080/", "application/json", strings.NewReader("{}"))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Printf("response: %T\n", resp)
	fmt.Printf("response.Body: %T\n", resp.Body)

	data := make([]byte, 8)
	for {
		readN, err := resp.Body.Read(data)
		if readN > 0 {
			fmt.Print(string(data[:readN]))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("done")

}
func server() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		flusher, ok := writer.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}

		for i := 0; i < 16; i++ {
			fmt.Fprintf(writer, "chunk [%02d]: %v\n", i, time.Now())
			flusher.Flush()
			time.Sleep(time.Second)
		}
	})

	http.ListenAndServe(":9080", nil)

}

func restyGET() {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"pageNum": "1",
			"limit":   "20",
			"sort":    "name",
			"order":   "asc",
			"random":  strconv.FormatInt(time.Now().Unix(), 10),
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("").
		Get("http://localhost:8085/userList")
	if err != nil {
		fmt.Println("err", err)
		return
	}
	pretty.Logln(string(resp.Body()))
}

type User struct {
	Name string `json:"name"`
}

func restyPost() {
	var u User
	client := resty.New()
	_, err := client.R().
		SetBody(map[string]interface{}{"username": "testuser", "password": "testpass"}).
		SetResult(&u). // or SetResult(AuthSuccess{}).
		Post("http://localhost:8085/login")
	if err != nil {
		fmt.Println("err", err)
	}
}
func uploadFileWithSharding() {

	// Create a Resty Client
	file, err := os.OpenFile("/Users/demg/Documents/ff.mp4", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	buf := make([]byte, 4096*1024*3)

	var i int32
	uid := cast.ToString(time.Now().UnixNano())
	res := map[string]interface{}{}
	for {
		i++
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read err", err)
			return
		}
		if err == io.EOF || n == 0 {
			fmt.Println("end")
			break
		}
		client := resty.New()
		_, err = client.R().
			SetFileReader("file", "file-"+cast.ToString(i), bytes.NewReader(buf[:n])).
			SetFormData(map[string]string{
				"fid": uid,
			}).
			SetResult(&res).
			Post("http://localhost:8085/upload")
		if err != nil {
			fmt.Println("upload file failed", err)
			return
		}

	}
	client := resty.New()
	_, err = client.R().
		SetFormData(map[string]string{
			"finish": "1",
			"fid":    uid,
		}).
		SetResult(&res).
		Post("http://localhost:8085/upload")
}

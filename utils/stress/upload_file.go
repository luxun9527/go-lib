package main

import (
	"bytes"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
)

func main() {
	flag.String("url", "", "url")
	flag.Parse()
	client := resty.New()
	//post 方法
	h := gin.H{}

	/*
		 setbody()
		      SetBody([]byte(`{"username":"testuser", "password":"testpass"}`)).
		      SetBody(`{"username":"testuser", "password":"testpass"}`).
		      SetBody(map[string]interface{}{"username": "testuser", "password": "testpass"}).
						User{
						Username: "jeeva@myjeeva.com",
						Password: "welcome2resty",
				})
	*/
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"username":"testuser", "password":"testpass"}`).
		SetResult(&h). // or SetResult(AuthSuccess{}).
		Post("http://localhost:9999/p1")
	if err != nil {
		log.Printf("Send Post Requst failed %v ", err)
		return
	}
	log.Printf("[post] body =%v result=%v", string(resp.Body()), h)

	//上传文件
	fs, err := os.ReadFile("./test.txt")
	if err != nil {
		log.Printf("open file failed %v", err)
		return
	}
	resp, err = client.R().
		SetFileReader("file", "test-img.png", bytes.NewReader(fs)).
		SetFormData(map[string]string{
			"first_name": "Jeevanandam",
			"last_name":  "M",
		}).
		SetResult(&h).
		Post("http://localhost:9999/upload")
	if err != nil {
		log.Printf("Send Upload Requst failed %v ", err)
		return
	}
	log.Printf("[POST] Upload body =%v result=%v", string(resp.Body()), h)
}

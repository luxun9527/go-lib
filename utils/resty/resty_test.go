package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	// Create a Resty Client
	client := resty.New().SetRetryCount(4)

	resp, err := client.R().
		EnableTrace().
		Get("https://www.baidu.com")

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}

// get post基本用法
func TestBaseGetPostBase(t *testing.T) {
	client := resty.New()
	// get方法
	h := gin.H{}
	resp, err := client.R().
		//SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more").
		SetQueryParams(map[string]string{
			"page_no": "1",
			"limit":   "20",
		}).SetHeaders(map[string]string{"token": "BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F"}).
		SetResult(&h).
		Get("http://localhost:9999/get")
	if err != nil {
		log.Printf("Send GET Requst failed %v ", err)
		return
	}
	log.Printf("[GET] body =%v result=%v", string(resp.Body()), h)

	//post 方法
	h = gin.H{}

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
	resp, err = client.R().
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
	//设置代理
	resp, err = client.SetProxy("http://127.0.0.1:7890").R().Get("http://localhost:9999/get")
	if err != nil {
		log.Printf("Test proxy body =%v result=%v", string(resp.Body()), h)
		return
	}
	log.Printf("[GET proxy] data %v", string(resp.Body()))
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())


}
func TestRetry(t *testing.T) {
	client := resty.New()
	client = client.AddRetryCondition(func(resp *resty.Response, err error) bool {
		return resp.StatusCode() == 200
	})
	// backoff to increase retry intervals after each attempt. 重试
	resp, err := client. // Set retry count to non zero to enable retries
		SetRetryCount(5).SetDebug(true).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(1 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(9 * time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		R().Get("http://localhost:9999/get")
	if err!=nil{
		log.Println(err)
	}
	log.Printf("[GET proxy] data %v", string(resp.Body()))
	select {

	}
}

func TestServer(t *testing.T) {
	engine := gin.New()
	engine.GET("/get", func(c *gin.Context) {
		c.Query("")

		log.Printf("limit=%v page_no=%v token=%v", c.Query("limit"), c.Query("page_no"), c.GetHeader("token"))
		c.JSON(200, gin.H{"code": "200", "msg": "success", "data": struct{}{}})

	})
	engine.POST("/p1", func(c *gin.Context) {
		h := gin.H{}
		if err := c.ShouldBindJSON(&h); err != nil {
			log.Println(err)
			return
		}
		log.Println(h)
		c.JSON(200, gin.H{"code": "200", "msg": "success", "data": struct{}{}})

	})
	engine.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(500, "上传图片出错")
		}
		// c.JSON(200, gin.H{"message": file.Header.Context})
		if err := c.SaveUploadedFile(file, file.Filename); err != nil {
			log.Printf("SaveUploadedFile failed %v", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"filename": file.Filename})
	})
	engine.Run(":9999")
}

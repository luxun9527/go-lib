package http

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":8888", nil)
}

// 处理 /upload  逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"sync"
)

type Query struct {
	PageNum int32  `form:"pageNum"`
	Limit   int32  `form:"limit"`
	Sort    string `form:"sort"`
	Order   string `form:"order"`
	Random  int32  `form:"random"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Recovery())
	//绑定get参数
	route.GET("/userList", QueryUserList)
	uploadHandler := NewUploadHandler()
	route.POST("/upload", uploadHandler.UploadFile)
	route.GET("/test", func(c *gin.Context) {
		c.JSON(200, "ok")
	})

	route.Run(":8085")
}

func QueryUserList(c *gin.Context) {
	var data Query
	if err := c.ShouldBindQuery(&data); err != nil {
		fmt.Println("err", err)
		return
	}
	c.JSON(200, data)
}

type UploadHandler struct {
	FileInfo map[string]*os.File
	lock     sync.RWMutex
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		FileInfo: map[string]*os.File{},
		lock:     sync.RWMutex{},
	}
}

func (u *UploadHandler) UploadFile(c *gin.Context) {

	fid := c.PostForm("fid")
	isFinish := c.PostForm("finish")
	if isFinish != "" {
		u.lock.Lock()
		if fs, ok := u.FileInfo[fid]; ok {
			fs.Close()
			delete(u.FileInfo, fid)
		}
		u.lock.Unlock()
		fmt.Printf("task finish fid = %v", fid)
	} else {
		u.lock.RLock()
		s, ok := u.FileInfo[fid]
		u.lock.RUnlock()
		if !ok {
			file, err := os.OpenFile("/Users/demg/personProject/go-lib/gin/ff.mp4", os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
			s = file
			u.lock.Lock()
			u.FileInfo[fid] = file
			u.lock.Unlock()
		}
		fd, err := c.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("debug 1 filename = %v\n", fd.Filename)
		f, err := fd.Open()
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		buf := make([]byte, 4096*1024*3)
		for {
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("read err", err)
				return
			}
			if err == io.EOF || n == 0 {
				break
			}
			if _, err := s.Write(buf[:n]); err != nil {
				fmt.Println("err", err)
				return
			}
		}

	}

	c.JSON(200, gin.H{"success": true})
}

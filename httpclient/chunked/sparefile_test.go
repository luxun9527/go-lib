package sparsefile

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"io"
	"log"
	url2 "net/url"
	"os"
	"testing"
)

func GenerateSpareFile() {
	distFs, err := os.OpenFile("test", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	defer distFs.Close()
	for i := 0; i < 100; i++ {
		r := cast.ToString(i)
		distFs.Seek(10000, io.SeekCurrent)
		distFs.Write([]byte(r))
	}

}

//go test -v sparefile_test.go sparsefile.go -test.run TestServer
func TestServer(t *testing.T) {
	r := gin.New()
	r.POST("/copy", func(c *gin.Context) {
		if err := ReceiveSparseFile(c, "testdist", func(i int64) {
			//log.Println("progress", i)
		}); err != nil {
			t.Error("err", err)
			return
		}
	})
	r.Run(":9094")
}

//go test -v sparefile_test.go sparsefile.go -test.run TestClient
func TestClient(t *testing.T) {
	GenerateSpareFile()
	file, err := os.Open("test")
	if err != nil {
		t.Fatal(err)
	}
	sparseFile := NewSparseFile(context.Background(), file)
	r2 := &url2.URL{
		Scheme: "http",
		Host:   "localhost:9094",
		Path:   "/copy",
	}
	err = sparseFile.CopyToRemote(r2, nil)
	if err != nil {
		t.Log("err", err)
	}

}

package utils

import (
	"context"
	"github.com/spf13/cast"
	"io"
	"log"
	"os"
	"testing"
)

func generateSpareFile() {
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

//go test -v sparsefile_test.go sparsefile.go -test.run TestCopy
func TestCopy(t *testing.T) {
	//	generateSpareFile()
	file, err := os.OpenFile("dist", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Log("err", err)
		return
	}

	if err := Copy(context.Background(), "/home/feng/iscsi-spare-vdisk-thin1668561994", file); err != nil {
		t.Log("err", err)
	}

}

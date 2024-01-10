package standard

import (
	"flag"
	"io/fs"
	"log"
	"path/filepath"
	"testing"
)

func TestWalk(t *testing.T) {
	filepath.Walk("S:\\go-lib\\standard", func(path string, info fs.FileInfo, err error) error {
		println(filepath.Ext(info.Name()))
		println(info.Name())
		return nil
	})
}
func TestReadDir(t *testing.T) {

	s := flag.String("GOPATH","","")
	flag.Parse()
	log.Println(*s)
}

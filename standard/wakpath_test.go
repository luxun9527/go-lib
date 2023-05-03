package standard

import (
	"fmt"
	"io/fs"
	"io/ioutil"
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
	files, err := ioutil.ReadDir("S:\\go-lib\\standard")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}
}

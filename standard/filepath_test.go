package standard

import (
	"log"
	"os"
	"path/filepath"

	"testing"
)

func TestFilePath(t *testing.T) {
	allFiles, err := getAllFiles("E:\\demoproject\\go-lib\\standard\\reflectdemo")
	if err != nil {
		log.Printf("get all files error:%v", err)
	}
	log.Printf("all files:%v", allFiles)
}

// 获取文件夹下所有文件
func getAllFiles(dir string) ([]string, error) {
	var files []string

	// 使用 filepath.Walk 遍历文件夹
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Printf("info %v", info.Name())
		// 如果是文件，则添加到结果列表中
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

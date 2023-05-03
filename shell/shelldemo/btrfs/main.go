package main

import (
	"io"
	"log"
	"os"
)

func main() {
	fd, err := os.OpenFile("/Volume3/src/test", os.O_RDWR, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	fileInfo, _ := fd.Stat()
	size := fileInfo.Size()
	//修改10%
	fd.Seek(io.SeekStart, 0)
	log.Println("size", size)
	var count int32
	for i := int64(0); i < size; i++ {
		if i%10 == 0 {
			fd.WriteAt([]byte("a"), i)
			count++
			if count == 3000000 {
				log.Println("finish")
				break
			}
		}

	}

}

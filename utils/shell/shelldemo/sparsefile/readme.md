用来拷贝稀疏文件如，iscsi lun 文件。

```go

package main

import (
	"log"

	"gitlab.local/golibrary/sparefile"
)

func main() {
	//	generateSpareFile()
	file, err := os.OpenFile("dist", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println("err", err)
		return
	}

	if err := sparefile.Copy(context.Background(), "/home/feng/iscsi-spare-vdisk-thin1668561994", file); err != nil {
		log.Println("err", err)
	}

}

```






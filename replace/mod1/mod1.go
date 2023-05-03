package testReplace

//在导出的时候go mod的包名会被替换
import "log"

func Export() {
	log.Println("aaa")
}

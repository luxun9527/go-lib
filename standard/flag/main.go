package main

import (
	"flag"
	"log"
)

func main() {
	s := flag.String("name", "", "用户名")
	//	go run main.go    --name=zhangsan -name=zhangsan   --name zhangsan -name zhangsan  2023/03/23 20:09:34 zhangsan
	var f float64
	flag.Float64Var(&f, "age", 1, "")
	var b bool
	//	go run main.go -f 2023/03/23 20:15:06 true
	flag.BoolVar(&b, "f", false, "")
	flag.Parse()
	log.Println(*s)
	log.Println(f)
	log.Println(b)
}

package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	parseEnv()
}

func baseUsage() {
	s := flag.String("name", "", "用户名")
	//	go run viper_test.go    --name=zhangsan -name=zhangsan   --name zhangsan -name zhangsan  2023/03/23 20:09:34 zhangsan
	var f float64
	flag.Float64Var(&f, "age", 1, "")
	var b bool
	//	go run viper_test.go -f 2023/03/23 20:15:06 true
	flag.BoolVar(&b, "f", false, "")
	flag.Parse()
	log.Println(*s)
	log.Println(f)
	log.Println(b)
}
func parseEnv(){
	s := flag.String("OS", "", "用户名")
	flag.Parse()
	log.Println(*s)
}
func flagSet() {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.Int64("age", 25, "")
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		log.Println(err)
	}
	flag.CommandLine.Int64("name", 25, "")
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		log.Println(err)
	}
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		log.Println(f)
	})
	flag.CommandLine.Visit(func(f *flag.Flag) {

	})

}

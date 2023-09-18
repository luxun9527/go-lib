package main

import (
	"fmt"
	"os"
	"strconv"
)

func  main()  {
	for idx, args := range os.Args {
		fmt.Println("参数" + strconv.Itoa(idx) + ":", args)
	}
}

//go run main.go 1 3 -X ?
//参数0: C:\Users\ADMINI~1\AppData\Local\Temp\go-build3741688355\b001\exe\main.exe
//参数1: 1
//参数2: 3
//参数3: -X
//参数4: ?

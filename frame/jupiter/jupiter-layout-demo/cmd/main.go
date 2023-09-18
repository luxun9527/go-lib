package main

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/engine"
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/model"
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/service"
	"log"
)

func main() {
	eng := engine.NewEngine()
	eng.RegisterHooks(hooks.Stage_AfterStop, func()  {
        fmt.Println("exit jupiter app ...")

      })

    model.Init()
    service.Init()
    if err := eng.Run(); err != nil {
    	log.Fatal(err)
    }
}


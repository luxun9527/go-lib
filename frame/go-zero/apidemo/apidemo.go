package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/rest"
	"go-lib/frame/go-zero/apidemo/internal/config"
	"go-lib/frame/go-zero/apidemo/internal/handler"
	"go-lib/frame/go-zero/apidemo/internal/svc"
	"log"
	"net/http"
	"runtime/debug"
)

var configFile = flag.String("f", "frame/go-zero/apidemo/etc/apidemo-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	_ = proc.AddWrapUpListener(func() {
		log.Println("testsetset")
	})
	logx.DisableStat()
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	server.Use(RecoverHandler)
	sseHandler := handler.NewSseHandler()

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/sse",
		Handler: sseHandler.Serve,
	}, rest.WithSSE())

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

type H map[string]interface{}

// RecoverHandler returns a middleware that recovers if panic happens.
func RecoverHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if result := recover(); result != nil {
				log.Printf("result %v stack string(%v)", result, string(debug.Stack()))
				data, _ := json.Marshal(H{"code": 1, "msg": "内部错误"})
				w.Write(data)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

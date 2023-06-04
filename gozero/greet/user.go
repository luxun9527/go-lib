package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"

	"greet/internal/config"
	"greet/internal/handler"
	"greet/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.DisableStat()
	logx.WithCallerSkip(1)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	//fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

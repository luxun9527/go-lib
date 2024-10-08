package engine

import (
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/worker/xcron"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/grpc/greeter"
	"go-lib/frame/jupiter/jupiter-layout-demo/internal/app/handler"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"time"
)

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		xgo.ParallelWithError(
			eng.serveGRPC,
			eng.serveHTTP,
			eng.startJobs,
		),
	); err != nil {
		xlog.Panic("startup engine", xlog.Any("err", err))
	}
	return eng
}

func (eng *Engine) serveHTTP() error {
	server := xecho.StdConfig("http").MustBuild()
	server.GET("/jupiter", func(ctx echo.Context) error {
		return ctx.JSON(200, "welcome to jupiter")
	})
	// Specify routing group
	group := server.Group("/api")
	group.GET("/user/:id", handler.GetUser)

	//support proxy for http to grpc controller
	g := greeter.Greeter{}
	group2 := server.Group("/grpc")
	group2.GET("/get", xecho.GRPCProxyWrapper(g.SayHello))
	group2.POST("/post", xecho.GRPCProxyWrapper(g.SayHello))
	return eng.Serve(server)
}

func (eng *Engine) serveGRPC() error {
	server := xgrpc.StdConfig("grpc").MustBuild()
	helloworld.RegisterGreeterServer(server.Server, new(greeter.Greeter))
	return eng.Serve(server)
}

func (eng *Engine) startJobs() error {
	cron := xcron.StdConfig("demo").Build()
	cron.Schedule(xcron.Every(time.Second*10), xcron.FuncJob(eng.execJob))
	return eng.Schedule(cron)
}

func (eng *Engine) execJob() error {
	xlog.Info("exec job", xlog.String("info", "print info"))
	xlog.Warn("exec job", xlog.String("warn", "print warning"))
	return nil
}

package initializer

import (
	"context"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/luxun9527/zlog"
	"go-lib/example/pkg/i18n"
	"go-lib/example/pkg/xgin"
	"go-lib/example/server/accountApi/config"
	"go-lib/example/server/accountApi/global"
	"go-lib/example/server/accountApi/rpcClient"
)

func Init(confPath string) {
	//初始化配置
	global.Config = config.InitConfig(confPath)
	//初始化日志
	zlog.InitDefaultLogger(&global.Config.Logger)
	//初始化grpc客户端
	client, err := global.Config.RpcClient.EtcdConf.NewEtcdClient()
	if err != nil {
		zlog.Panicf("etcd client init failed, err:%v", err)
	}
	if err := rpcClient.InitEtcdRpcClients(context.Background(), client, global.Config.RpcClient.TargetConfList); err != nil {
		zlog.Panicf("rpc client init failed, err:%v", err)
	}
	//加载多语言文件
	translator, err := i18n.NewTranslatorFormFile(global.Config.Lang.Path)
	if err != nil {
		zlog.Panicf("i18n init failed, err:%v", err)
	}
	i18n.SetDefaultTranslator(translator)

	//设置gin参数校验失败翻译
	v, _ := xgin.NewValidateTranslator(binding.Validator.Engine().(*validator.Validate))
	xgin.SetDefaultValidateTranslator(v)
}

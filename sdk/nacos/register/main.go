package main

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

type CustomLogger  struct {}

func (l *CustomLogger) Infof(format string, v ...interface{}) { 
    log.Printf(format, v...)
}


func (l *CustomLogger) Warnf(format string, v ...interface{}) {
	log.Printf(format, v...)
}


func (l *CustomLogger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}


func (l *CustomLogger) Debugf(format string, v ...interface{}) {
	log.Printf(format, v...)
}


func (l *CustomLogger) Info(v ...interface{}) {
	log.Println(v...)
}

func (l *CustomLogger) Warn(v ...interface{}) {
	log.Println(v...)
}
func (l *CustomLogger) Error(v ...interface{}) {
	log.Println(v...)
}
func (l *CustomLogger) Debug(v ...interface{}) {
	log.Println(v...)
}

func main() {
	//create ServerConfig

	sc := []constant.ServerConfig{
		*constant.NewServerConfig("192.168.11.185", 8848, constant.WithContextPath("/nacos")),
	}
	cc := constant.ClientConfig{
		TimeoutMs:            0,
		ListenInterval:       0,
		BeatInterval:         0,
		NamespaceId:          "2177cd45-6846-4b56-8dcc-3df3f7d6280a",
		AppName:              "",
		AppKey:               "",
		Endpoint:             "",
		RegionId:             "",
		AccessKey:            "",
		SecretKey:            "",
		OpenKMS:              false,
		KMSVersion:           "",
		KMSv3Config:          nil,
		CacheDir:             "",
		DisableUseSnapShot:   false,
		UpdateThreadNum:      0,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: false,
		Username:             "nacos",
		Password:             "nacos",
		LogDir:               "D:\\project\\go-lib\\sdk\\nacos\\register\\tmp",
		LogLevel:             "debug",
		ContextPath:          "",
		AppendToStdout:       false,
		LogSampling:          nil,
		LogRollingConfig:     nil,
		TLSCfg:               constant.TLSConfig{},
		AsyncUpdateService:   false,
		EndpointContextPath:  "",
		EndpointQueryParams:  "",
		ClusterName:          "",
	}
	//create ClientConfig

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil{
		log.Printf("new config client failed, err:%v", err)
	}
	logger.SetLogger(&CustomLogger{})
	config, err := client.GetConfig(vo.ConfigParam{
		DataId:           "test",
		Group:            "test",
		Content:          "",
		Tag:              "",
		AppName:          "",
		BetaIps:          "",
		CasMd5:           "",
		Type:             "",
		SrcUser:          "",
		EncryptedDataKey: "",
		KmsKeyId:         "",
		UsageType:        "",
		OnChange:         nil,
	})
	if err !=nil{
		log.Panicf("err:%v", err)
	}
	log.Printf("config:%v", config)
}

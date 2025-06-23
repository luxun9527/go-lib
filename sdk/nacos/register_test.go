package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"testing"
)

type CustomLogger struct{}

func (l *CustomLogger) Infof(format string, v ...interface{}) {
	log.Println("Infof")
	log.Printf(format, v...)
}

func (l *CustomLogger) Warnf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *CustomLogger) Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *CustomLogger) Debugf(format string, v ...interface{}) {
	log.Println("Debugf")
	log.Printf(format, v...)
}

func (l *CustomLogger) Info(v ...interface{}) {
	log.Println("info")
	log.Println(v...)
}

func (l *CustomLogger) Warn(v ...interface{}) {
	log.Println("warn")
	log.Println(v...)
}
func (l *CustomLogger) Error(v ...interface{}) {
	log.Println("error")
	log.Println(v...)
}
func (l *CustomLogger) Debug(v ...interface{}) {
	log.Println("Debug")
	log.Println(v...)
}

func TestRegister(t *testing.T) {
	//create ServerConfig

	sc := []constant.ServerConfig{
		*constant.NewServerConfig("192.168.31.100", 8848, constant.WithContextPath("/nacos")),
	}
	cc := constant.ClientConfig{
		TimeoutMs:            0,
		ListenInterval:       0,
		BeatInterval:         0,
		NamespaceId:          "05932850-81ea-422b-873c-809a963627d7",
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
		LogDir:               "./tmp",
		LogLevel:             "warn",
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

	if err != nil {
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
	if err != nil {
		log.Panicf("err:%v", err)
	}
	log.Printf("config:%v", config)
}

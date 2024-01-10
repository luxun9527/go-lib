package main

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

func main() {
	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("192.168.11.185", 8848, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(3000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("D:\\project\\go-lib\\sdk\\nacos\\register\\tmp"),
		constant.WithCacheDir("D:\\project\\go-lib\\sdk\\nacos\\register\\tmp"),
	//	constant.WithLogLevel("debug"),
		constant.WithUsername("nacos"),
		constant.WithPassword("nacos"),
	)

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

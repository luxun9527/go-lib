package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"testing"
)

func TestPublishConfig(t *testing.T) {
	// 定义 Nacos 服务器地址
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("192.168.31.100", 8848, constant.WithContextPath("/nacos")),
	}

	// 定义客户端配置
	cc := constant.ClientConfig{
		NamespaceId:         "05932850-81ea-422b-873c-809a963627d7",
		TimeoutMs:           5000,
		ListenInterval:      30000,
		Username:            "nacos",
		Password:            "nacos",
		LogDir:              "./tmp",
		LogLevel:            "debug",
		NotLoadCacheAtStart: true,
	}

	// 创建配置客户端
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		t.Fatalf("Failed to create config client: %v", err)
	}

	// 上传配置
	dataId := "test1"
	group := "test"
	content := "key=value\nanotherKey=anotherValue"
	success, err := client.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
	if err != nil || !success {
		t.Errorf("Failed to publish config: %v", err)
	}
	log.Printf("Successfully published config: DataId=%s, Group=%s", dataId, group)
}

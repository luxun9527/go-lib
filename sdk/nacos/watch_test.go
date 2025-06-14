package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"sync"
	"testing"
	"time"
)

func TestWatchConfig(t *testing.T) {
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

	// 等待客户端就绪（可优化为检查状态）
	time.Sleep(3 * time.Second)

	// 监听配置
	dataId := "test"
	group := "test"

	var wg sync.WaitGroup
	wg.Add(1)

	d, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			log.Printf("Config changed! Namespace: %s, Group: %s, DataId: %s, New Content: %s", namespace, group, dataId, data)
			wg.Done()
		},
	})
	if err != nil {
		t.Errorf("Failed to listen config: %v", err)
		return
	}

	log.Printf("Listening for config changes: DataId=%s, Group=%s d=%v", dataId, group, d)

	time.Sleep(5 * time.Hour)
}

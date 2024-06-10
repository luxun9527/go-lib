package viper

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"log"
	"net/http"
	"testing"
	"time"
)

// refer https://github.com/spf13/viper
func TestExample(t *testing.T) {
	viper.SetConfigFile("./config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	ip := viper.Get("servers.alpha.ip")
	fmt.Printf("ip = %v \n", ip)
}
func TestSetConfig(t *testing.T) {
	viper.SetConfigFile("./config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	viper.Set("name", "zhangsan")
	if err := viper.WriteConfig(); err != nil {
		log.Println(err)
	}
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName("new")
	//SafeWriteConfig 要 SetConfigType 和AddConfigPath 不能修改已经存在的配置文件必须新建。要覆盖可以使用WriteConfig
	if err := viper.SafeWriteConfig(); err != nil {
		log.Println(err)
	}
	if err := viper.WriteConfigAs("config.toml"); err != nil {
		log.Println(err)
	}
}

func TestUnmarshal(t *testing.T) {
	//mapstrcture tag 指定配置文件中的变量
	type Config struct {
		Favorite []string `mapstructure:"fav"`
		Name     string
		AgeAge   int32

		Age string `mapstructure:"age"`
	}
	type C2 struct {
		Name string
	}
	viper.SetConfigFile("./config.toml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return
	}
	var c2 C2
	if err := viper.Unmarshal(&c2); err != nil {
		return
	}
	sub := viper.Sub("person")
	//config := make(map[string]interface{})
	var config Config
	if err := sub.Unmarshal(&config); err != nil {
		fmt.Println("err = ", err)
		return
	}
	fmt.Printf("config %+v \n", config)
}
func TestWatchFile(t *testing.T) {
	viper.SetConfigFile("./config.toml") // 指定配置文件
	err := viper.ReadInConfig()          // 读取配置信息
	if err != nil {                      // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("change")
	})
	// 监控配置文件变化
	viper.WatchConfig()
	r := gin.Default()
	// 访问/version的返回值会随配置文件的变化而变化
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("file.version"))
	})

	if err := r.Run(
		fmt.Sprintf(":%d", viper.GetInt("gin.port"))); err != nil {
		log.Panic(err)
	}
}
func TestRemoteConfig(t *testing.T) {
	v := viper.New()
	err := v.AddRemoteProvider("etcd3", "http://192.168.2.159:2379", "language/zh-CN")
	if err != nil {
		log.Println(err)
		return
	}
	v.SetConfigType("yaml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	if err := v.ReadRemoteConfig(); err != nil {
		log.Println(err)
		return
	}
	result := v.Get("110000")
	log.Println(result)

	go func() {
		for {
			//相当于一直主动查询
			time.Sleep(time.Second * 3)
			if err := v.WatchRemoteConfig(); err != nil {
				log.Println(err)
				return
			}
			result := v.Get("110000")
			log.Printf("tragger %v result %v", "WatchRemoteConfig", result)
		}
	}()

	go func() {
		for {
			//主动推送当有修改的时候推送。
			if err := v.WatchRemoteConfigOnChannel(); err != nil {
				log.Println(err)
				return
			}
			result := v.Get("110000")
			log.Printf("tragger %v result %v", "WatchRemoteConfigOnChannel", result)
		}
	}()

}

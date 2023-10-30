package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	//example1()
	exampleUnmarshal()
	//exampleWatchFile()
}
func example1() {
	viper.SetConfigFile("E:\\demoproject\\go-lib\\utils\\viper\\config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	ip := viper.Get("servers.alpha.ip")
	fmt.Printf("ip = %v \n", ip)
}

func exampleUnmarshal() {
	//mapstrcture tag 指定配置文件中的变量
	type Config struct {
		Favorite []string `mapstructure:"fav"`
		Name     string
		AgeAge   int32

		Age string `mapstructure:"age"`
	}
	viper.SetConfigFile("./utils/viper/config.toml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
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
func exampleWatchFile() {
	viper.SetConfigFile("./config.toml") // 指定配置文件
	err := viper.ReadInConfig()          // 读取配置信息
	if err != nil {                      // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

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

package config

import (
	"github.com/luxun9527/zlog"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConf  `mapstructure:"server"`
	Logger    zlog.Config `mapstructure:"logger"`
	RpcClient RpcClient   `mapstructure:"rpcClient"`
	Lang      struct {
		Path string
	}
}

func InitConfig(path string) *Config {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		zlog.Panicf("viper.ReadInConfig failed, err:%v", err)
	}
	var c Config
	if err := viper.Unmarshal(&c, viper.DecodeHook(zlog.StringToLogLevelHookFunc())); err != nil {
		zlog.Panicf("viper.Unmarshal failed, err:%v", err)
	}

	return &c

}

package config

import "github.com/luxun9527/zlog"

type Config struct {
	Server ServerConf  `mapstructure:"server"`
	Logger zlog.Config `mapstructure:"logger"`
}

func InitConfig() {

}

package config

import (
	"github.com/luxun9527/zlog"
	"github.com/spf13/viper"
	"go-lib/example/pkg/xgorm"
	"go-lib/example/pkg/xjwt"
	"go-lib/example/pkg/xredis"
)

type Config struct {
	Server    ServerConf       `mapstructure:"server"`
	Logger    zlog.Config      `mapstructure:"logger"`
	GormConf  xgorm.GormConf   `mapstructure:"gromConf"`
	JwtConf   xjwt.JwtConf     `mapstructure:"jwtConf"`
	RedisConf xredis.RedisConf `mapstructure:"redisConf"`
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
	if c.GormConf.Logger == nil {
		loggerConf := c.Logger
		c.GormConf.Logger = &loggerConf
	}
	c.Server.RegisterConf.Port = c.Server.Port

	return &c

}

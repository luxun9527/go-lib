package config

import "go-lib/example/pkg/xetcd"

type ServerConf struct {
	Port         int32                  `mapstructure:"port"`
	RegisterConf xetcd.EtcdRegisterConf `mapstructure:"registerConf"`
}

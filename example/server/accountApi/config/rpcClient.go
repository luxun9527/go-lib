package config

import "go-lib/example/pkg/xetcd"

type RpcClient struct {
	EtcdConf       xetcd.EtcdConf
	TargetConfList []*TargetConf
}

type TargetConf struct {
	Key     string
	TimeOut int64
}

package xetcd

import (
	"github.com/luxun9527/zlog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdConf struct {
	Hosts []string
}

func (c *EtcdConf) NewEtcdClient() (*clientv3.Client, error) {

	config := clientv3.Config{
		Endpoints: c.Hosts, DialTimeout: time.Second * time.Duration(5),
		Logger: zlog.DefaultLogger,
	}
	return clientv3.New(config)
}

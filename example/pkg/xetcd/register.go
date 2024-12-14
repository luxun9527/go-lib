package xetcd

import (
	"context"
	"github.com/luxun9527/zlog"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/netx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type EtcdRegisterConf struct {
	EtcdConf
	Key   string
	Value string
	Port  int32
}

func Register(conf EtcdRegisterConf) {
	go func() {
		cli, err := conf.EtcdConf.NewEtcdClient()
		if err != nil {
			zlog.Panicf("etcd new client err: %v", err)
		}
		manager, err := endpoints.NewManager(cli, conf.Key)
		if err != nil {
			zlog.Panicf("etcd new manager err: %v", err)
		}
		//设置租约时间
		resp, err := cli.Grant(context.Background(), 5)
		if err != nil {
			zlog.Panicf("etcd grant err: %v", err)
		}
		if conf.Value == "" {
			conf.Value = netx.InternalIp() + ":" + cast.ToString(conf.Port)
		}
		if err := manager.AddEndpoint(context.Background(), conf.Key+"/"+cast.ToString(int64(resp.ID)), endpoints.Endpoint{Addr: conf.Value}, clientv3.WithLease(resp.ID)); err != nil {
			zlog.Panicf("etcd add endpoint err: %v", err)
		}
		c, err := cli.KeepAlive(context.Background(), resp.ID)
		if err != nil {
			zlog.Panicf("etcd keepalive err: %v", err)
		}
		zlog.Infof("etcd register success,key: %v,value: %v", conf.Key, conf.Value)
		for {
			select {
			case _, ok := <-c:
				if !ok {
					zlog.Errorf("etcd keepalive failed,please check etcd key %v existed", conf.Key)
					return
				}
			}
		}

	}()

}

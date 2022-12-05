package service

import (
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/spf13/viper"

	"helper/common/errno"
	"helper/common/logger"
)

type Reply struct {
	errno.CodeErr
	Data interface{} `json:"data"`
}

func StartService(server *server.Server) {
	address := viper.GetString("etcd.pub_host")
	etcdUrl := viper.GetString("etcd.url")

	//etcd 注册中心
	rplugin := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + address,
		EtcdServers:    []string{etcdUrl},
		BasePath:       viper.GetString("etcd.pub_base_dir"),
		Metrics:        metrics.NewRegistry(),
		Services:       make([]string, 0),
		UpdateInterval: 30 * time.Second,
	}
	err := rplugin.Start()
	if err != nil {
		logger.Instance.WithField("code", errno.ErrService).Panicf("error start rplugin: %v", err)
	}
	server.Plugins.Add(rplugin)
}

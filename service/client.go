package service

import (
	"context"

	etcdC "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/protocol"
	"github.com/spf13/viper"

	"helper/common/errno"
	"helper/common/logger"
)

type ServiceClient struct {
	Name string
}

func NewServiceClient(name string) *ServiceClient {
	return &ServiceClient{
		Name: name,
	}
}

func (s *ServiceClient) CallMapItf(method string, args interface{}) (map[string]interface{}, errno.CodeErr) {
	data, err := s.Call(method, args)
	return data.(map[string]interface{}), err
}

func (s *ServiceClient) Call(method string, args interface{}) (interface{}, errno.CodeErr) {
	reply, err := send(s.Name, method, args)
	if err != nil {
		logger.Instance.WithField("code", errno.ErrService).Panicf("error %v call: %v.", s.Name, err)
	}
	return reply.Data, reply.CodeErr
}

func send(service, method string, args interface{}) (*Reply, error) {
	d, err := etcdC.NewEtcdV3Discovery(viper.GetString("etcd.base_dir"), service, []string{viper.GetString("etcd.url")}, true, nil)
	if err != nil {
		return nil, err
	}

	opt := client.DefaultOption
	opt.SerializeType = protocol.JSON
	xClient := client.NewXClient(service, client.Failtry, client.RandomSelect, d, opt)
	defer func() {
		d.Close()
		err = xClient.Close()
		if err != nil {
			logger.Instance.WithField("code", errno.ErrService).Panicf("rpc client close err")
		}
	}()
	log.SetLogger(logger.Instance)

	reply := &Reply{}
	err = xClient.Call(context.Background(), method, args, reply)
	logger.Instance.Debugf("method: %v; args:%+v; err:%+v; reply:%#v.", service+"."+method, args, err, *reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

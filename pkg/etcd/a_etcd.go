package etcd

import (
	"errors"
	"github.com/golang/glog"
	"sync"
)

var (
	cli  *EtcdClient
	rw   *ReadWriter
	once sync.Once
)

func Connect(cfg *Config, clsName string) (err error) {
	glog.Infof("Setting up etcd.")
	once.Do(func() {
		cli = NewEtcdClient(cfg, clsName)
		if cli != nil {
			rw = NewEtcdReadWriter(cli)
		}
	})

	if cli == nil {
		return errors.New("failed to initialize etcd")
	}

	return nil
}

func Close() {
	glog.Infof("Closing etcd.")
}

func GetClsReadWriter() *ReadWriter {
	return rw
}

func GetEtcdCli() *EtcdClient {
	return cli
}

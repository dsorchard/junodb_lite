package config

import "junodb_lite/pkg/c_etcd"

type Config struct {
	EtcdEnabled bool
	Etcd        etcd.Config
	ClusterName string
	Outbound    interface{}
}

var (
	Conf = Config{}
)

package config

import "junodb_lite/pkg/etcd"

type Config struct {
	EtcdEnabled bool
	Etcd        etcd.Config
	ClusterName string
}

var (
	Conf = Config{}
)

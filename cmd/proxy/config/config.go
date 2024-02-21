package config

import "junodb_lite/pkg/b_etcd"

type Config struct {
	EtcdEnabled bool
	Etcd        etcd.Config
	ClusterName string
}

var (
	Conf = Config{}
)

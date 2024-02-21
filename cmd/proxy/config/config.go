package config

import "junodb_lite/pkg/a_etcd"

type Config struct {
	EtcdEnabled bool
	Etcd        etcd.Config
	ClusterName string
}

var (
	Conf = Config{}
)

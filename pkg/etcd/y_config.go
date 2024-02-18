package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Config struct {
	clientv3.Config
	RequestTimeout     time.Duration
	MaxConnectAttempts int
	MaxConnectBackoff  int
	CacheDir           string
	CacheName          string
	EtcdKeyPrefix      string
}

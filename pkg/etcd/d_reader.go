package etcd

// Implements cluster.IReader
type Reader struct {
	etcdcli *EtcdClient
}

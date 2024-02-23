package etcd

import clientv3 "go.etcd.io/etcd/client/v3"

// List of transactional operations
type OpList []clientv3.Op

func (op *OpList) AddDeleteWithRange(beginKey string, endKey string) {
	*op = append(*op, clientv3.OpDelete(beginKey, clientv3.WithRange(endKey)))
}

func (op *OpList) AddDeleteWithPrefix(key string) {
	*op = append(*op, clientv3.OpDelete(key, clientv3.WithPrefix()))
}

func (op *OpList) AddPut(key string, val string) {
	*op = append(*op, clientv3.OpPut(key, val))
}

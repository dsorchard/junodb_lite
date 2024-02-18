package etcd

import "fmt"

type IKVWriter interface {
	PutValue(key string, value string) (err error)
	DeleteKeyWithPrefix(key string, isPrefix bool) (err error)
	PutValuesWithTxn(op OpList) (err error)
}

type BaseWriter struct {
	kvwriter IKVWriter
}

type EtcdWriter struct {
	BaseWriter
	etcdcli *EtcdClient
}

type StdoutWriter struct {
	BaseWriter
	keyPrefix string
}

func (w *EtcdWriter) PutValue(key string, value string) (err error) {
	return w.etcdcli.PutValue(key, value)
}

func (w *EtcdWriter) DeleteKeyWithPrefix(key string, isPrefix bool) (err error) {
	return w.etcdcli.DeleteKeyWithPrefix(key, isPrefix)
}

func (w *EtcdWriter) PutValuesWithTxn(op OpList) (err error) {
	return w.etcdcli.PutValuesWithTxn(op)
}

func (w *StdoutWriter) PutValue(key string, value string) (err error) {
	fmt.Printf("%s%s=%s\n", w.keyPrefix, key, value)
	return nil
}

func (w *StdoutWriter) DeleteKeyWithPrefix(key string, isPrefix bool) (err error) {
	fmt.Printf("delete: key=%s%s isPrefix=%v\n", w.keyPrefix, key, isPrefix)
	return nil
}

func (w *StdoutWriter) PutValuesWithTxn(op OpList) (err error) {
	if len(op) == 0 {
		return nil
	}

	fmt.Printf("===txn begin:\n")
	for i := 0; i < len(op); i++ {
		if op[i].IsDelete() {
			fmt.Printf("delete: beginKey=%s%s\n", w.keyPrefix, string(op[i].KeyBytes()))
			endKey := op[i].RangeBytes()
			if endKey != nil {
				fmt.Printf("          endKey=%s%s\n", w.keyPrefix, string(endKey))
			}
		} else {
			fmt.Printf("%s%s=%s\n", w.keyPrefix,
				string(op[i].KeyBytes()), string(op[i].ValueBytes()))
		}
	}
	fmt.Printf("ops_count=%d\n", len(op))
	fmt.Printf("===txn end\n")

	return nil
}

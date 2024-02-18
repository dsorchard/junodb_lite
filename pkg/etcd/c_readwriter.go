package etcd

type ReadWriter struct {
	Reader
	EtcdWriter
}

func NewEtcdReadWriter(cli *EtcdClient) *ReadWriter {
	rw := &ReadWriter{
		Reader: Reader{
			etcdcli: cli,
		},
		EtcdWriter: EtcdWriter{
			etcdcli: cli,
		},
	}

	// self for polymorphism to work
	rw.kvwriter = rw
	return rw
}

type EtcdReadStdoutWriter struct {
	Reader
	StdoutWriter
}

func NewEtcdReadStdoutWriter(cli *EtcdClient, clusterName string) *EtcdReadStdoutWriter {
	rw := &EtcdReadStdoutWriter{
		Reader: Reader{
			etcdcli: cli,
		},
		StdoutWriter: StdoutWriter{
			keyPrefix: cli.config.EtcdKeyPrefix + clusterName + TagCompDelimiter,
		},
	}

	// self for polymorphism to work
	rw.kvwriter = rw
	return rw
}

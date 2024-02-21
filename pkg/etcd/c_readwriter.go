package etcd

type ReadWriter struct {
	Reader
	Writer
}

func NewEtcdReadWriter(cli *Client) *ReadWriter {
	rw := &ReadWriter{
		Reader: Reader{
			etcdcli: cli,
		},
		EtcdWriter: Writer{
			etcdcli: cli,
		},
	}

	// self for polymorphism to work
	rw.kvwriter = rw
	return rw
}

type ReadStdoutWriter struct {
	Reader
	StdoutWriter
}

func NewEtcdReadStdoutWriter(cli *Client, clusterName string) *ReadStdoutWriter {
	rw := &ReadStdoutWriter{
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

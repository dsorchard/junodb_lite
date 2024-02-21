package etcd

var chanForProxy chan int

func WatchForProxy() chan int {
	if chanForProxy != nil {
		return chanForProxy
	}

	chanForProxy = make(chan int, 2)
	if cli == nil {
		close(chanForProxy)
	} else {
		go cli.WatchEvents(TagVersion, chanForProxy)
	}

	return chanForProxy
}

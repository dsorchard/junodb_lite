package watcher

import (
	"context"
	"github.com/golang/glog"
	cluster "junodb_lite/pkg/b_cluster"
	etcd "junodb_lite/pkg/c_etcd"
	"sync"
)

type Watcher struct {
	clustername string
	etcdcli     *etcd.Client
	etcdcfg     *etcd.Config
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

var (
	theWatcher  *Watcher
	markdownobj *cluster.ZoneMarkDown
)

func Initialize(args ...interface{}) (err error) {
	return nil
}
func Finalize() {
	if theWatcher != nil {
		theWatcher.Stop()
	}
}
func (w *Watcher) Stop() {
	glog.Infof("stop watcher")
	w.cancel()
	w.wg.Wait()
}

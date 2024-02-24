package replication

import (
	"fmt"
	"github.com/golang/glog"
	repconfig "junodb_lite/cmd/group1/a_proxy/b_config"
	proto "junodb_lite/pkg/ac_proto"
	io "junodb_lite/pkg/y_conn_mgr"
	util "junodb_lite/pkg/y_util"
	"sync"
)

var (
	TheReplicator *Replicator
	initOnce      sync.Once
	enabled       bool = false
)

type (
	Replicator struct {
		conf       *repconfig.Config
		processors []*replicationProcessorT
	}
	replicationProcessorT struct {
		io.OutboundProcessor
		reqCtxCreator repReqCtxCreatorI
		specNsMap     map[string]bool
		byPassLTM     bool
	}
	repReqCtxCreatorI interface {
		newRequestContext(recExpirationTime uint32, msg *proto.RawMessage, reqCh chan io.IRequestContext,
			dropCnt *util.AtomicShareCounter, errCnt *util.AtomicShareCounter) io.IRequestContext
		newKeepAliveRequestContext() io.IRequestContext
	}
)

func (r *Replicator) Shutdown() {
	for _, processor := range r.processors {
		processor.Shutdown()
	}

	for _, processor := range r.processors {
		processor.WaitShutdown()
	}
}

func Initialize(args ...interface{}) (err error) {
	sz := len(args)
	if sz == 0 {
		err = fmt.Errorf("replication config expected")
		glog.Error(err)
		return
	}
	conf, ok := args[0].(*repconfig.Config)
	if !ok {
		err = fmt.Errorf("wrong argument type")
		glog.Error(err)
		return
	}
	err = Init(conf)
	return
}

func Init(conf *repconfig.Config) error {
	return nil
}

func Finalize() {
	if enabled && TheReplicator != nil {
		TheReplicator.Shutdown()
	}
}

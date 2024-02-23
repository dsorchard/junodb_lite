package cli

import (
	proto "junodb_lite/pkg/a_proto"
	"junodb_lite/pkg/io"
	"sync"
	"time"
)

type IOError struct {
	Err error
}

func (e *IOError) Error() string {
	return "IOError: " + e.Err.Error()
}

var (
	kMaxRequestChanBufferSize = 1024

	kCalTxnType    = "JUNO_CLIENT"
	kCalSslTxnType = "JUNO_SSL_CLIENT"
)

type Processor struct {
	server     io.ServiceEndpoint
	sourceName string

	connectTimeout     time.Duration
	requestTimeout     time.Duration
	connRecycleTimeout time.Duration

	chDone     chan bool
	chProcDone <-chan bool
	chRequest  chan *RequestContext
	startOnce  sync.Once
}

func NewProcessor(
	server io.ServiceEndpoint,
	sourceName string,
	connectTimeout time.Duration,
	requestTimeout time.Duration,
	connRecycleTimeout time.Duration) *Processor {

	c := &Processor{
		server:             server,
		sourceName:         sourceName,
		connectTimeout:     connectTimeout,
		requestTimeout:     requestTimeout,
		connRecycleTimeout: connRecycleTimeout,
		chDone:             make(chan bool),
		chRequest:          make(chan *RequestContext, kMaxRequestChanBufferSize),
	}
	return c
}

func (c *Processor) Start() {
	c.startOnce.Do(func() {
		//c.chProcDone = StartRequestProcessor(
		//	c.server, c.sourceName, c.connectTimeout, c.requestTimeout, c.connRecycleTimeout, c.chDone, c.chRequest)
	})
}

func (c *Processor) Close() {

}

func (c *Processor) ProcessRequest(request *proto.OperationalMessage) (*proto.OperationalMessage, error) {
	return nil, nil
}

func (c *Processor) ProcessBatchRequests(requests []*proto.OperationalMessage) (responses []*proto.OperationalMessage, err error) {
	return nil, nil
}

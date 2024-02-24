package service

import (
	"github.com/golang/glog"
	io "junodb_lite/pkg/y_conn_mgr"
	"os"
	"sync"
)

type Service struct {
	listeners      []io.IListener
	wg             sync.WaitGroup
	chDone         chan bool
	chSuspend      chan bool // true: suspend, false: resume
	suspended      bool
	config         Config
	requestHandler io.IRequestHandler
	inShutdown     int32 ///TODO to be renamed
	acceptLimiter  ILimiter
	Zoneid         int
}

func (s *Service) Run() {

}

func NewWithLimiterAndListenFd(cfg Config, reqHandler io.IRequestHandler, limiter ILimiter, fds ...*os.File) (service *Service) {
	service = NewWithListenFd(cfg, reqHandler, fds...)
	service.acceptLimiter = limiter
	return
}

func NewWithListenFd(cfg Config, reqHandler io.IRequestHandler, fds ...*os.File) (service *Service) {
	var listeners []io.IListener
	cfgListeners := cfg.Listener
	if len(cfgListeners) != len(fds) {
		glog.Fatal("number of listener config not match number of FDs")
	}

	for i, fd := range fds {

		ln, err := io.NewListenerWithFd(cfgListeners[i], cfg.GetIoConfig(&cfgListeners[i]), fd, reqHandler)
		if err == nil {
			listeners = append(listeners, ln)
		} else {
			glog.Warning("Cannot Listen on ", fd.Fd())
		}
	}
	return New(cfg, reqHandler, listeners...)

}
func New(config Config, reqHandler io.IRequestHandler, listeners ...io.IListener) (service *Service) {
	service = &Service{
		listeners:      listeners,
		chDone:         make(chan bool),
		chSuspend:      make(chan bool, 1),
		config:         config,
		requestHandler: reqHandler,
		Zoneid:         -1,
	}
	return
}

type SuspendFunc func(b bool)

func NewService(cfg Config, reqHandler io.IRequestHandler) (*Service, SuspendFunc) {
	//cfg.SetDefaultIfNotDefined()

	cfgListeners := cfg.Listener

	var listeners []io.IListener
	for _, lsnrCfg := range cfgListeners {

		ln, err := io.NewListener(lsnrCfg, cfg.GetIoConfig(&lsnrCfg), reqHandler)
		if err == nil {
			listeners = append(listeners, ln)
		} else {
			glog.Warningf("Cannot Listen on %s, err=%s", lsnrCfg.Addr, err.Error())
		}
	}
	if len(listeners) == 0 {
		glog.Fatal("No listener created")
	}

	svc := New(cfg, reqHandler, listeners...)
	f := func(b bool) {
		svc.chSuspend <- b
	}
	return svc, f
}

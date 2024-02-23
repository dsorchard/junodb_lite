package service

import (
	"fmt"
	"github.com/golang/glog"
	"net"
	"os"
	"os/exec"
	"strings"
)

type (
	ChildInfo struct {
		Id  int
		Cmd *exec.Cmd
	}

	httpMonitoringT struct {
		listeners []net.Listener
		lsnrFiles []*os.File
	}
	ServerManager struct {
		monitoring *httpMonitoringT

		listeners    []net.Listener
		files        []*os.File
		cmdPath      string
		cmdArgs      []string
		numChildren  int
		pidMap       map[int]ChildInfo
		doneCh       chan struct{}
		deadCh       chan int
		stopping     bool
		cloudEnabled bool
	}
)

func NewServerManager(num int, path string, args []string, cfg Config, httpMonAddr string, cloudEnabled bool) *ServerManager {
	s := &ServerManager{
		cmdPath:      path,
		cmdArgs:      args,
		numChildren:  num,
		pidMap:       make(map[int]ChildInfo),
		doneCh:       make(chan struct{}),
		deadCh:       make(chan int, num),
		stopping:     false,
		cloudEnabled: cloudEnabled,
	}
	if len(httpMonAddr) != 0 { ///TODO validate addr?
		s.monitoring = &httpMonitoringT{
			listeners: make([]net.Listener, s.numChildren),
			lsnrFiles: make([]*os.File, s.numChildren),
		}
		for i := 0; i < s.numChildren; i++ {
			if ln, err := net.Listen("tcp", "127.0.0.1:0"); err == nil {
				s.monitoring.listeners[i] = ln
				s.monitoring.lsnrFiles[i], _ = ln.(*net.TCPListener).File()
				glog.Infof("monitoring listener (%s) created for child %d", ln.Addr(), i)

			} else {
				glog.Errorf("fail to create listeners for monitoring")
			}
		}
	}

	cfgListeners := cfg.Listener
	numListeners := len(cfgListeners)

	s.listeners = make([]net.Listener, numListeners)

	s.files = make([]*os.File, numListeners)

	for i, lsnrCfg := range cfgListeners {

		ln, err := net.Listen("tcp", lsnrCfg.Addr)
		if err != nil {
			glog.Exitf("Cannot Listen on %s", lsnrCfg.Addr)
		}
		s.listeners[i] = ln
		file, _ := ln.(*net.TCPListener).File()
		s.files[i] = file
	}

	return s
}

func (s *ServerManager) Run() {
	s.handleSignals()

	hostName, err := os.Hostname()
	spawn := true
	if err == nil && s.cloudEnabled {
		shutdownList := fmt.Sprintf(" %s ", os.Getenv("SHUTDOWN_LIST"))
		name := fmt.Sprintf(" %s ", hostName)
		glog.Infof("host=%s shutdownList=%s", hostName, shutdownList)
		if strings.Contains(shutdownList, " all ") || strings.Contains(shutdownList, name) {
			glog.Infof("Skip starting workers on:%s", shutdownList)
			spawn = false
		}
	}

	if spawn {
		s.spawnChildren()
	}

Loop:
	for {
		select {
		//case pid := <-s.deadCh:
		//s.handleDeadChild(pid)
		case <-s.doneCh:
			//s.shutdown()
			break Loop
		}
	}
}

func (s *ServerManager) handleSignals() {

}

func (s *ServerManager) spawnChildren() {
	for i := 0; i < s.numChildren; i++ {
		s.spawnOneChild(i)
	}
	s.spawnMonitoringChild()
}

func (s *ServerManager) spawnMonitoringChild() {
	if s.monitoring != nil {
		var addrs []string
		for _, ln := range s.monitoring.listeners {
			addrs = append(addrs, ln.Addr().String())
		}
		if len(addrs) == 0 {
			return
		}

		param := fmt.Sprintf("-worker-monitoring-addresses=%s", strings.Join(addrs, ","))
		var args []string = []string{
			"monitor",
			"-child",
			param,
		}
		args = append(args, s.cmdArgs...)

		cmd := exec.Command(s.cmdPath, args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = s.files

		err := cmd.Start()
		if err != nil {
			glog.Fatalf("Failed to launch child process, error: %v", err)
		}

		// save the cmd for later
		s.pidMap[cmd.Process.Pid] = ChildInfo{-1, cmd}
		glog.Infof("%s %s", s.cmdPath, strings.Join(args, " "))
	}
}

func (s *ServerManager) spawnOneChild(i int) {
	param := fmt.Sprintf("-worker-id=%d", i)
	var args []string = []string{
		"worker",
		param,
		"-child",
	}
	args = append(args, s.cmdArgs...)

	cmd := exec.Command(s.cmdPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if s.monitoring != nil {
		cmd.ExtraFiles = append(s.files, s.monitoring.lsnrFiles[i])
	} else {
		cmd.ExtraFiles = s.files
	}

	err := cmd.Start()
	if err != nil {
		glog.Fatalf("Failed to launch child process, error: %v", err)
	}

	// save the cmd for later
	s.pidMap[cmd.Process.Pid] = ChildInfo{i, cmd}
	glog.Infof("%s %s", s.cmdPath, strings.Join(args, " "))
}

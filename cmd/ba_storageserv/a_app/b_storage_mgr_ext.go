package app

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type (
	ConnectInfo struct {
		Listener     string
		ZoneId       int
		MachineIndex int
	}
	ChildInfo struct {
		Id  int
		Cmd *exec.Cmd
	}
	ServerManager struct {
		//monitoring *httpMonitoringT

		connectInfo        []ConnectInfo
		pidFileName        string
		cmdPath            string
		cmdArgs            []string
		numChildren        int
		pidMap             map[int]ChildInfo
		doneCh             chan struct{}
		deadCh             chan int
		stopping           bool
		procCreateItvlBase int
		procCreateItvlMax  int
		lruCacheSizeInMB   int
		dbScanPort         int
		cloudEnabled       bool
	}
)

func NewServerManager(num int, pidFileName string, path string, args []string,
	connInfo []ConnectInfo, httpMonAddr string, dbScanPort int, cloudEnabled bool) *ServerManager {
	s := &ServerManager{
		connectInfo:        connInfo,
		pidFileName:        pidFileName,
		cmdPath:            path,
		cmdArgs:            args,
		numChildren:        num,
		pidMap:             make(map[int]ChildInfo),
		doneCh:             make(chan struct{}),
		deadCh:             make(chan int, num),
		stopping:           false,
		procCreateItvlBase: 100,
		procCreateItvlMax:  20000,
		// For rocksdb.  min(10% * mem / num_of_db_instances, 3072)
		//lruCacheSizeInMB: int(math.Min(float64(util.GetTotalMemMB()/(10*num)), 3072)),
		dbScanPort:   dbScanPort,
		cloudEnabled: cloudEnabled,
	}
	if len(httpMonAddr) != 0 { ///TODO validate addr?
		for i := 0; i < s.numChildren; i++ {
			if ln, err := net.Listen("tcp", "127.0.0.1:0"); err == nil {
				//s.monitoring.listeners[i] = ln
				//s.monitoring.lsnrFiles[i], _ = ln.(*net.TCPListener).File()
				glog.Infof("monitoring listener (%s) created for worker %d", ln.Addr(), i)

			} else {
				glog.Errorf("fail to create listeners for monitoring")
			}
		}
	}

	return s
}

func (s *ServerManager) Run() {
	pidFile := s.pidFileName

	if data, err := os.ReadFile(pidFile); err == nil {
		if pid, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			if process, err := os.FindProcess(pid); err == nil {
				if err := process.Signal(syscall.Signal(0)); err == nil {
					fmt.Fprintf(os.Stderr, "process pid: %d in %s is still running\n", pid, pidFile)
					///TODO check if it is storageserv process
					os.Exit(-1)
				}
			}
		}
	}

	if s.dbScanPort > 0 {
		cmdPath := fmt.Sprintf("%s/%s", filepath.Dir(s.cmdPath), "dbscanserv")
		_, err := os.Stat(cmdPath)
		if errors.Is(err, fs.ErrNotExist) {
			glog.Exitf("missing executable file: dbscanserv.")
		}
	}

	os.WriteFile(pidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644)
	defer os.Remove(pidFile)
	//defer shmstats.Finalize()

	//if err := shmstats.InitForManager(s.numChildren); err != nil {
	//	glog.Error(err.Error())
	//	return
	//}

	//s.handleSignals()

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
		case pid := <-s.deadCh:
			s.handleDeadChild(pid)
		case <-s.doneCh:
			//s.shutdown()
			break Loop
		}
	}
}
func (s *ServerManager) spawnChildren() {
	for i := 0; i < s.numChildren; i++ {
		s.spawnOneChild(i)
	}
	if s.dbScanPort > 0 {
		s.spawnDbScanChild()
	}
}

func (s *ServerManager) spawnDbScanChild() {

	cmdPath := fmt.Sprintf("%s/%s", filepath.Dir(s.cmdPath), "dbscanserv")

	var argConfig string
	for _, val := range s.cmdArgs {
		if strings.Index(val, "-config") == 0 || strings.Index(val, "-c") == 0 {
			argConfig = val
			break
		}
	}
	glog.Infof("%s %s", s.cmdPath, argConfig)
	cmd := exec.Command(cmdPath, argConfig)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		glog.Fatalf("Failed to launch child process, error: %v", err)
	}

	// save the cmd for later
	s.pidMap[cmd.Process.Pid] = ChildInfo{-2, cmd}

}
func (s *ServerManager) spawnOneChild(i int) {

}
func (s *ServerManager) handleDeadChild(pid int) {
}

package initmgr

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	initializers initEntriesT
)

type entryT struct {
	initializer  IInitializer
	weight       int
	args         []interface{}
	initOnce     sync.Once
	finalizeOnce sync.Once
}

type initEntriesT []entryT

type IInitializer interface {
	Name() string
	Initialize(args ...interface{}) error
	Finalize()
}

func (rs initEntriesT) Len() int {
	return len(rs)
}

func (rs initEntriesT) Less(i, j int) bool {
	return rs[i].weight < rs[j].weight
}

func (rs initEntriesT) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func Init() {
	signal.Ignore(syscall.SIGPIPE, syscall.SIGURG)

	sigDoneCh := make(chan bool)
	sigCh := make(chan os.Signal, 10)
	signal.Notify(sigCh)
	defer func() {
		close(sigDoneCh)
		signal.Stop(sigCh)
	}()

	go func(sigCh chan os.Signal) {
		fmt.Fprintf(os.Stderr, "... signal handling started\n")
		for {
			select {
			case <-sigDoneCh:
				return
			case sig := <-sigCh:
				fmt.Fprintf(os.Stderr, "... signal %d (%s) received\n", sig, sig)

				switch sig {
				case syscall.SIGTERM, syscall.SIGINT, syscall.SIGSEGV, syscall.SIGBUS:
					os.Stderr.Sync()
					os.Exit(0)
				}
			}
		}
	}(sigCh)
	sort.Sort(initializers)
	var err error

	for i, _ := range initializers {
		initializers[i].initOnce.Do(func() {
			name := initializers[i].initializer.Name()
			if err = initializers[i].initializer.Initialize(initializers[i].args...); err == nil {
				fmt.Fprintf(os.Stderr, "... [ok]   initmgr.initialize %s\n", name)
			} else {
				fmt.Fprintf(os.Stderr, "... [fail] initmgr.initialize %s\t (error: %s)\n", name, err.Error())
			}
		})
		if err != nil {
			tm := time.Second * 5
			tm = tm + time.Duration(rand.Int63n(int64(tm)))
			time.Sleep(tm)
			fmt.Fprintf(os.Stderr, "\n... Initialization FAILURE. Exit in %s ...\n\n", tm.String())

			finalizeBackwardsFrom(i - 1)
			os.Stderr.Sync()
			os.Exit(255)
		}
	}
}

func finalizeBackwardsFrom(i int) {
	for ; i >= 0; i-- {
		initializers[i].finalizeOnce.Do(func() {
			name := initializers[i].initializer.Name()
			fmt.Fprintf(os.Stderr, "... initmgr.finalize %s\n", name)
			initializers[i].initializer.Finalize()
		})
	}
}
func Finalize() {
	finalizeBackwardsFrom(len(initializers) - 1)
}

func Register(rc IInitializer, args ...interface{}) {
	RegisterWithWeight(rc, len(initializers), args...)
}

func RegisterWithFuncs(initializeFunc func(args ...interface{}) error, finalizeFunc func(), args ...interface{}) {
	Register(NewInitializer(initializeFunc, finalizeFunc), args...)
}

func RegisterWithWeight(rc IInitializer, weight int, args ...interface{}) {
	initializers = append(initializers, entryT{initializer: rc, weight: weight, args: args})
}

type Initializer struct {
	name           string
	InitializeFunc func(args ...interface{}) error
	FinalizeFunc   func()
}

func (i *Initializer) Name() string {
	return i.name
}

func (i *Initializer) Initialize(args ...interface{}) (err error) {
	if i.InitializeFunc != nil {
		if err = i.InitializeFunc(args...); err != nil {
			return
		}
	}
	return
}

func (i *Initializer) Finalize() {
	if i.FinalizeFunc != nil {
		i.FinalizeFunc()
	}
}

func NewInitializer(initializeFunc func(args ...interface{}) error, finalizeFunc func()) IInitializer {
	name := runtime.FuncForPC(reflect.ValueOf(initializeFunc).Pointer()).Name()
	i := strings.LastIndex(name, ".")
	if i == -1 {
		name = "unknown package"
	} else {
		name = name[0:i]
	}
	return &Initializer{name, initializeFunc, finalizeFunc}
}

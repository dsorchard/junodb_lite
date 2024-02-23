package util

import "syscall"

func GetNumOpenFDs() (n int) {
	// alternatives on Unix/Linux:
	// * /proc/<pid>/fd
	// * lsof
	var rlim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim); err == nil {
		for i := 0; i < int(rlim.Cur); i++ {
			var stat syscall.Stat_t
			if e := syscall.Fstat(i, &stat); e == nil {
				n++
			}
		}
	}
	return
}

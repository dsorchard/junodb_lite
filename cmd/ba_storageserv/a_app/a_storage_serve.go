package app

import (
	"fmt"
	cmd "junodb_lite/pkg/d_cmd"
	"os"
)

func init() {
	var (
		cmdManager Manager
		cmdWorker  Worker
	)
	cmdManager.Init("manager", "start as storage server manager")
	cmdWorker.Init("worker", "start as storage worker")
	cmd.Register(&cmdManager)
	cmd.Register(&cmdWorker)
}

func Main() {
	indexCommand := 1
	numArgs := len(os.Args)

	if indexCommand < numArgs {
		cmd := cmd.GetCommand(os.Args[indexCommand])
		if cmd != nil {
			if err := cmd.Parse(os.Args[indexCommand+1:]); err == nil {
				cmd.Exec()
			} else {
				fmt.Printf("* command '%s' failed. %s\n", cmd.GetName(), err)
			}
		} else {
			fmt.Printf("command '%s' not specified", os.Args[indexCommand])
			return
		}
	}
}

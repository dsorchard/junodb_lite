package app

import (
	"fmt"
	cmd "junodb_lite/pkg/d_cmd"
	initmgr "junodb_lite/pkg/e_initmgr"
)

func Main() {
	defer initmgr.Finalize()

	var (
		cmdManager Manager
		cmdWorker  Worker
	)
	cmdManager.Init("manager", "start as proxy server manager")
	cmdWorker.Init("worker", "start as proxy worker")
	cmd.Register(&cmdManager)
	cmd.Register(&cmdWorker)

	if command, args := cmd.ParseCommandLine(); command != nil {
		if err := command.Parse(args); err == nil {
			command.Exec()
		} else {
			fmt.Printf("* command '%s' failed. %s\n", command.GetName(), err)
		}
	} else {
		//execDefault()
	}
}

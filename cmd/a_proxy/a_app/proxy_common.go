package app

import (
	"errors"
	"fmt"
	"io/fs"
	cmd "junodb_lite/pkg/d_cmd"
	"os"
)

type (
	CmdProxyCommon struct {
		cmd.Command
		optConfigFile string
		optLogLevel   string
	}
)

func (c *CmdProxyCommon) Init(name string, desc string) {
	//c.Command.Init(name, desc)
	//c.StringOption(&c.optConfigFile, "c|config", "", "specify toml config file")
	//c.StringOption(&c.optLogLevel, "log-level", kDefaultLogLevel, "specify log level")
}

func (c *CmdProxyCommon) Parse(args []string) (err error) {
	if err = c.Command.Parse(args); err != nil {
		return
	}
	if len(c.optConfigFile) == 0 {
		fmt.Fprintf(os.Stderr, "\n\n*** missing config option ***\n\n")
		c.FlagSet.Usage()
		os.Exit(-1)
	}

	if _, err := os.Stat(c.optConfigFile); errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "\n\n***  config file \"%s\" not found ***\n\n", c.optConfigFile)
		os.Exit(-1)
	}

	return
}

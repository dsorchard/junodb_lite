package cmd

import "flag"

type (
	Option struct {
		flag.FlagSet
		optsDesc string
	}
)

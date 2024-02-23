package cmd

import (
	"fmt"
	"os"
)

var (
	commands           = make(map[string]ICommand)
	groups             = make(map[string]*Group)
	notGroupedCommands []ICommand
)

type (
	ICommand interface {
		GetName() string
		GetDesc() string //get short description
		GetSynopsis() string
		GetDetails() string
		GetOptionDesc() string
		GetExample() string
		AddExample(cmdExample string, desc string)
		AddDetails(txt string)
		Init(name string, desc string)
		Exec()
		Parse(args []string) error
		PrintUsage()
	}

	Command struct {
		Option
		name       string
		desc       string //short description. (one ine)
		synopsis   string
		details    string
		examples   string
		optVModule string
	}

	Group struct {
		cmds []ICommand
		name string
	}
)

func Register(c ICommand) bool {
	if register(c) {
		notGroupedCommands = append(notGroupedCommands, c)
		return true
	}
	return false
}
func register(c ICommand) bool {
	if _, found := commands[c.GetName()]; found {
		fmt.Printf("Command %s has been registered.", c.GetName())
		return false
	}
	commands[c.GetName()] = c
	return true
}

func ParseCommandLine() (cmd ICommand, args []string) {
	numArgs := len(os.Args)

	for i := 1; i < numArgs; i++ {
		arg := os.Args[i]
		if cmd == nil {
			cmd = GetCommand(arg)
			if cmd != nil {
				args = append(args, os.Args[i+1:]...)
				break
			}
		}
		args = append(args, arg)
	}
	return
}
func GetCommand(name string) ICommand {
	if cmd, ok := commands[name]; ok {
		return cmd
	}
	return nil
}

package app

type (
	Manager struct {
		CmdStorageCommon
		optNumChildren uint
		optIpAddress   string
		cmdArgs        []string
	}
)

func (s *Manager) GetName() string {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) GetDesc() string {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) GetSynopsis() string {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) GetDetails() string {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) GetOptionDesc() string {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) GetExample() string {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) AddExample(cmdExample string, desc string) {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) AddDetails(txt string) {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) Exec() {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) PrintUsage() {
	//TODO implement me
	panic("implement me")
}

func (s *Manager) Init(name string, desc string) {
}

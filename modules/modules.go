package modules

type Modules struct {
	Hub map[string]*BaseModule
	DonePipe chan struct{}
}

type BaseModule struct {
	Module

	ID uint8
	Status uint32
}

type Module interface {
	Configure(*Modules, ...interface{}) (Module, error)
	//	Unconfigure()

	//Start() error
	//Stop() error
}

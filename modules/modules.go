package modules

type Modules struct {
	Ids map[string]uint8
	Hub map[uint8]*BaseModule
}

type BaseModule struct {
	Module

	ID uint8
	Status uint32
	Error_ch chan error
}

type Module interface {
	Configure(*Modules, ...interface{}) (Module, error)
	//	Unconfigure()

	//Start() error
	//Stop() error
}

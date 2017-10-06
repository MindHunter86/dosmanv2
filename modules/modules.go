package modules

import "sync"

type Modules struct {
	// BaseModule address "storage":
	Hub map[string]*BaseModule

	// Modules global control:
	DonePipe chan struct{}
	WaitGroup sync.WaitGroup
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

package modules

import (
	"sync"
	config "mailru/rooster22/system/config"
	"github.com/rs/zerolog"
)

type Modules struct {
	// BaseModule address "storage":
	Hub map[string]*BaseModule

	// Modules global resources:
	Logger *zerolog.Logger
	Config *config.SysConfig

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
	Bootstrap() error
}

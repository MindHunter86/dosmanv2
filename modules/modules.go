package modules

import (
	"sync"
	config "mailru/rooster22/system/config"
	"github.com/rs/zerolog"
)

const (
	StatusReady = uint32(iota) // call when module has been configured (pre start state)
	StatusRunning // call when module bootstraped successfully
	StatusStopping // call when we want stop module
	StatusFailed // call, when module has failed on* configure or bootstrap ("failed on" or "failed in"?)
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

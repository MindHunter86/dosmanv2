package modules

import (
	"sync"
	config "mailru/rooster22/system/config"
	"github.com/rs/zerolog"
)

const (
	StatusReady = uint8(iota) // call when module has been configured (pre start state)
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
	ErrorPipe chan *ModuleError
	WaitGroup sync.WaitGroup
}

type Module interface {
	Configure(*Modules, ...interface{}) (Module, error)
	Bootstrap() error
}

type BaseModule struct {
	Module
	status uint8
}

type ModuleError struct {
	ModName string
	ErrLevel string

	E error
}


// ModuleError API:
func (self *ModuleError) Error() error { return self.E }


// BaseModules API:
func (self *BaseModule) GetModuleStatus() uint8 {
	return self.status
}
func (self *BaseModule) SetModuleStatus(status uint8) {
	self.status = status
}

package modules

import (
	"sync"
	config "mh00appserver/system/config"
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
	e error
	modName string
}


// ModuleError API:
func (self *ModuleError) Error() error { return self.e }
func (self *ModuleError) SetError(e error) *ModuleError {
	self.e = e
	return self
}
func (self *ModuleError) ModuleName() string { return self.modName }
func (self *ModuleError) SetModuleName(modName string) *ModuleError {
	self.modName = modName
	return self
}


// BaseModules API:
func (self *BaseModule) GetModuleStatus() uint8 {
	return self.status
}
func (self *BaseModule) SetModuleStatus(status uint8) {
	self.status = status
}

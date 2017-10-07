package system

import (
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"mailru/rooster22/modules"
	"mailru/rooster22/modules/http"
	"mailru/rooster22/modules/mysql"
	config "mailru/rooster22/system/config"

	"github.com/rs/zerolog"
)


// System structs:
type System struct {
	wg *sync.WaitGroup
	log *zerolog.Logger
	cfg *config.SysConfig

	mods *modules.Modules
}


// System API:
func (self *System) SetLogger(logger *zerolog.Logger) *System {
	self.log = logger
	return self
}

func (self *System) Configure() (*System, error) {
	var e error

	// parse configuration file:
	if self.cfg, e = new(config.SysConfig).Parse(); e != nil { return nil,e }

	// define new modulelist:
	self.mods = new(modules.Modules)
	self.mods.Hub = make(map[string]*modules.BaseModule)

	self.mods.Logger = self.log
	self.mods.DonePipe = make(chan struct{})
	if self.mods.Config, e = new(config.SysConfig).Parse(); e != nil { return nil,e }

	// modules loader:
	if e = self.preloadModule(new(http.HttpModule).Configure(self.mods, nil)); e != nil { return nil,e }
	if e = self.preloadModule(new(mysql.MysqlModule).Configure(self.mods, nil)); e != nil { return nil,e }

	return self,nil
}

func (self *System) Launch() error {
	for modName, modPointer := range self.mods.Hub { go self.moduleBootstrap(modName, modPointer)	}
	return nil
}

func (self *System) LaunchEvLoops() error {
	var kernSignal chan os.Signal = make(chan os.Signal)
	signal.Notify(kernSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

DSTR:
	for {
		select {
		case <-kernSignal:
			self.log.Warn().Msg("Syscall.SIG* has been detected! Closing application...")
			break DSTR
		}
	}

	return self.destroy()
}


// System internal methods:
func (self *System) moduleBootstrap(modName string, modPointer *modules.BaseModule) {
	self.log.Info().Msg("Bootstrapping module \""+modName+"\"...")

	self.mods.Hub[modName].Status = modules.StatusRunning // Check StatusReady ???
	self.mods.WaitGroup.Add(1)

	if e := modPointer.Bootstrap(); e != nil && self.mods.Hub[modName].Status != modules.StatusStopping {
		self.mods.WaitGroup.Done()
		self.mods.Hub[modName].Status = modules.StatusFailed
		self.log.Error().Err(e).Msg("Recieved error from module \""+modName+"\"! Status changed to StatusFailed.")
		return
	}

	self.mods.WaitGroup.Done()
	self.mods.Hub[modName].Status = modules.StatusReady
	self.log.Info().Msg("Module \""+modName+"\" has been successfully stopped and unloaded!")
}

func (self *System) preloadModule(modPointer modules.Module, modError error) error {
	// fail app if new module has an error:
	if modError != nil { return modError }

	// append new module to map:
	var modName reflect.Type = reflect.TypeOf(modPointer)

	self.mods.Hub[modName.Elem().Name()] = &modules.BaseModule{
		ID: uint8(len(self.mods.Hub)),
		Module: modPointer,
		Status: modules.StatusReady,
	}

	self.log.Debug().Str("module", modName.Elem().Name()).Msg("Module has been configured! Status changed to StatusReady.")
	return nil
}

func (self *System) destroy() error {
	// prepare:
	for modName, _ := range self.mods.Hub {
		self.mods.Hub[modName].Status = modules.StatusStopping
	}

	// closing:
	close(self.mods.DonePipe)
	self.mods.WaitGroup.Wait()

	return nil
}

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
	"mailru/rooster22/modules/telegram"
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
	if e = self.preloadModule(new(telegram.TelegramModule).Configure(self.mods, nil)); e != nil { return nil,e }

	return self,nil
}

func (self *System) Bootstrap() error {
	// define global error var for modError pipe:
	var e error
	var modErrorPipe chan *modules.ModuleError = make(chan *modules.ModuleError)

	// define kernel signal listener:
	var kernelSignal chan os.Signal = make(chan os.Signal)
	signal.Notify(kernelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	// bootstrap configured modules:
	for modName, modPointer := range self.mods.Hub {
		if self.mods.Hub[modName].GetModuleStatus() == modules.StatusReady {
			go self.moduleBootstrap(modName, modPointer, modErrorPipe)
		}
	}

	// start sys event loop:
LOOP:
	for {
		select {
		case <-kernelSignal:
			self.log.Warn().Msg("Syscall.SIG* has been detected! Closing application...")
			break LOOP
		case modError := <-modErrorPipe:
			e = modError.Error()
			self.log.Error().Str("MODULE", modError.ModuleName()).Err(e).Msg("CRITICAL ERROR!")
			break LOOP
		}
	}

	// TODO: Add buf for modErrorPipe. Check self.mods.WaitGroup
	// TODO: Check module order!!! 

	// stop and unload all modules:
	self.shutdown()

	// return nil or errors from modules (over mods.ErrorPipe):
	return e
}


// System internal methods:
func (self *System) preloadModule(modPointer modules.Module, e error) error {
	// fail app if new module has an error:
	if e != nil { return e }

	// append new module to map:
	var modName string = reflect.TypeOf(modPointer).Elem().Name()
	self.mods.Hub[modName] = &modules.BaseModule{ Module: modPointer }
	self.mods.Hub[modName].SetModuleStatus(modules.StatusReady)

	self.log.Debug().Str("module", modName).Msg("Module has been configured! Status changed to StatusReady.")
	return nil
}

func (self *System) moduleBootstrap(modName string, modPointer *modules.BaseModule, modError chan *modules.ModuleError) {
	self.log.Info().Msg("Bootstrapping module \""+modName+"\"...")

	self.mods.Hub[modName].SetModuleStatus(modules.StatusRunning)
	self.mods.WaitGroup.Add(1)

	if e := modPointer.Bootstrap(); e != nil && self.mods.Hub[modName].GetModuleStatus() != modules.StatusStopping {
		modError<- new(modules.ModuleError).SetModuleName(modName).SetError(e)
		self.mods.Hub[modName].SetModuleStatus(modules.StatusFailed)
		self.mods.WaitGroup.Done()
		return
	}

	self.mods.WaitGroup.Done()
	self.mods.Hub[modName].SetModuleStatus(modules.StatusReady)
	self.log.Info().Msg("Module \""+modName+"\" has been successfully stopped and unloaded!")
}

func (self *System) shutdown() {
	// Set status STOP for running modules:
	for modName, _ := range self.mods.Hub {
		if self.mods.Hub[modName].GetModuleStatus() == modules.StatusRunning {
			self.mods.Hub[modName].SetModuleStatus(modules.StatusStopping)
		}
	}

	// Close done pipe and wait modules unload:
	close(self.mods.DonePipe)
	self.mods.WaitGroup.Wait()
}

package system

import (
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"mailru/rooster22/modules"
	tgrm "mailru/rooster22/modules/telegram"
	config "mailru/rooster22/system/config"

	"github.com/rs/zerolog"
)


type System struct {
	wg *sync.WaitGroup
	log *zerolog.Logger
	cfg *config.SysConfig

	mods *modules.Modules
}

func (self *System) Configure() (*System, error) {
	var e error

	// parse configuration file:
	if self.cfg, e = new(config.SysConfig).Parse(); e != nil { return nil,e }

	// 
	self.mods = new(modules.Modules)
	self.mods.Ids = make(map[string]uint8)
	self.mods.Hub = make(map[uint8]*modules.BaseModule)

	// load modules:
	for { // error "catcher":
		if e = self.preloadModule(new(tgrm.TelegramBot).Configure(self.mods,nil)); e != nil { break }
		break
	}
	if e != nil {
		self.log.Debug()
		return nil,e
	}

	return self,nil
}
func (self *System) SetLogger(logger *zerolog.Logger) *System {
	self.log = logger
	return self
}
func (self *System) LaunchEvLoops() error {
	var kernSignal chan os.Signal = make(chan os.Signal)
	signal.Notify(kernSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

DSTR:
	for {
		select {
		case <-kernSignal:
			self.log.Debug().Msg("Kernel signal caught!")
			break DSTR
		}
	}

	return self.destroy()
}

func (self *System) preloadModule(modPointer modules.Module, modError error) error {
	// fail app if new module has an error:
	if modError != nil { return modError }

	// append new module to map:
	var modID uint8 = uint8(len(self.mods.Hub))
	self.mods.Hub[modID] = &modules.BaseModule{
		ID: modID,
		Module: modPointer,
	}

	// append new module id to map:
	modName := reflect.TypeOf(modPointer)
	self.mods.Ids[modName.Elem().Name()] = modID

	return nil
}
func (self *System) destroy() error {
	return nil
}

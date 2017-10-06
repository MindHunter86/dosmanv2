package system

import (
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"mailru/rooster22/modules"
	//tgrm "mailru/rooster22/modules/telegram"
	config "mailru/rooster22/system/config"
	"mailru/rooster22/modules/mysql"

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
	self.mods.Hub = make(map[string]*modules.BaseModule)
	self.mods.DonePipe = make(chan struct{})

	// load modules:
	for { // error "catcher":
		// temporary disabled module:
		//if e = self.preloadModule(new(tgrm.TelegramBot).Configure(self.mods,nil)); e != nil { break }
		if e = self.preloadModule(new(mysql.MysqlModule).Configure(self.mods, self.log, self.cfg)); e != nil { break }
		break
	}
	if e != nil {
		self.log.Error().Err(e).Msg("Preload module configuration hase been failed!")
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
	self.mods.WaitGroup.Add(1)
	var modName reflect.Type = reflect.TypeOf(modPointer)

	self.mods.Hub[modName.Elem().Name()] = &modules.BaseModule{
		ID: uint8(len(self.mods.Hub)),
		Module: modPointer,
	}

	return nil
}
func (self *System) destroy() error {
	close(self.mods.DonePipe)
	self.mods.WaitGroup.Wait()

	return nil
}

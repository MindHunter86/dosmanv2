package http

import (
	"reflect"

	"mailru/rooster22/modules"

	"github.com/rs/zerolog"
//	"github.com/buaazp/fasthttprouter"
)

type HttpModule struct {
	log zerolog.Logger

	modName string
	mods *modules.Modules
	donePipe chan struct{}
}

func (self *HttpModule) Configure(mods *modules.Modules, args ...interface{}) (modules.Module, error) {
	self.mods = mods
	self.modName = reflect.TypeOf(self).Elem().Name()

	// Set module name as prefix for logger:
	self.log = self.mods.Logger.With().Str("module", self.modName).Logger()

	go self.startCloseEventLoop()
	return self,nil
}
func (self *HttpModule) Start() error { return nil }
func (self *HttpModule) Stop() error { return nil }
func (self *HttpModule) Unconfigure() {}

func (self *HttpModule) startCloseEventLoop() {
	<-self.mods.DonePipe
	self.mods.WaitGroup.Done()
	self.log.Debug().Msg("donePipe closed!")
}

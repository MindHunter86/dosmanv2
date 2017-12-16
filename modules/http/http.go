package http

import "reflect"
import "mh00appserver/modules"
import "github.com/rs/zerolog"


// Module structs:
type HttpModule struct {
	log zerolog.Logger
	http *httpServer

	modName string
	mods *modules.Modules
	donePipe chan struct{}
}


// Module API:
func (self *HttpModule) Configure(mods *modules.Modules, args ...interface{}) (modules.Module, error) {
	var e error

	// get pointer for modulelist && get module title (reflect gets main struct title):
	self.mods = mods
	self.modName = reflect.TypeOf(self).Elem().Name()

	// Set module name as prefix for logger:
	self.log = self.mods.Logger.With().Str("module", self.modName).Logger()

	// define new httpServer:
	if self.http, e = new(httpServer).configure(self.mods.Config); e != nil { return nil,e }

	// start "exit signal catcher":
	go self.parentEventHandler()
	return self,nil
}
func (self *HttpModule) Bootstrap() error { return self.http.serve() }


// Module internal functions:
func (self *HttpModule) parentEventHandler () {
	// wait parent stop (done) signal:
	<-self.mods.DonePipe

	// close httpServer socket for fasthttp fatal error:
	self.http.socket.Close()
}

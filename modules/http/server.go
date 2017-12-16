package http

import "net"
import config "mh00appserver/system/config"
import "github.com/valyala/fasthttp"


type httpServer struct {
	router *httpRouter
	socket net.Listener
}


// Define and Configuration:
func (self *httpServer) configure(cfg *config.SysConfig) (*httpServer, error) {
	var e error

	// create network socket for future serving:
	if self.socket, e = net.Listen("tcp4", cfg.Base.Http.Listen); e != nil { return nil,e }
	// define http router for REST api:
	self.router = new(httpRouter).configure()

	return self,nil
}
// Simple serving by fasthttp:
func (self *httpServer) serve() error { return fasthttp.Serve(self.socket, self.router.Handler) }

package http

import "github.com/valyala/fasthttp"
import "github.com/buaazp/fasthttprouter"


type httpRouter struct {
	*fasthttprouter.Router
}


// Define and Configuration:
func (self *httpRouter) configure() (*httpRouter) {
	// define new router:
	self.Router = fasthttprouter.New()

	// configure httpapi methods:
	self.GET("/", self.index)

	return self
}


// Router handlers:
func (self *httpRouter) index(ctx *fasthttp.RequestCtx) {
	// ./router.go:26: cannot use ([]byte)("Hello world!") (type []byte) as type *bufio.Writer in argument to ctx.Response.Write
	// ctx.Response.Write([]byte("Hello world!"))
}

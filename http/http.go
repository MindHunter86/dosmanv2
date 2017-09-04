package http

import (
	"mailru/rooster22/config"

	"golang.org/x/net/context"
	"github.com/rs/zerolog"
	"github.com/buaazp/fasthttprouter"
)

type httpService struct {
	appLogger *zerolog.Logger
	appConfig *config.AppConfig

	httpRouter *fasthttprouter.Router
	httpController *httpController

	ctxPipeDone <-chan struct{}
}
func (self *httpService) ConfigureAndServe(ctx context.Context) error {
	self.ctxPipeDone = ctx.Done()
	self.appConfig = ctx.Value(config.CTX_APP_CONFIG).(*config.AppConfig)
	self.appLogger = ctx.Value(config.CTX_APP_LOGGER).(*zerolog.Logger)

	self.httpController = new(httpController).New(self)
	return nil
}

func (self *httpService) createHttpRouter() {
	self.httpRouter = fasthttprouter.New()
//	self.GET("/")
}

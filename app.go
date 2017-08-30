package main

import (
	"sync"
	"os"
	"os/signal"
	"syscall"

	"mailru/rooster22/config"
	"mailru/rooster22/telegram"

	"golang.org/x/net/context"
	"github.com/rs/zerolog"
)


type application struct {
	wg *sync.WaitGroup
	log *zerolog.Logger
	cfg *config.AppConfig

	tgbot *telegram.TelegramBot

	compCtxClose context.CancelFunc
}

// pub functions:
func (self *application) CreateAndConfigure() (*application, error) {
	var e error
	var ctx context.Context

	// parse configuration file:
	if self.cfg, e = new(config.AppConfig).Parse(); e != nil { return nil,e }

	// create context and fill it with required data:
	ctx, self.compCtxClose = context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, config.CTX_APP_LOGGER, self.log)
	ctx = context.WithValue(ctx, config.CTX_APP_CONFIG, self.cfg)

	// components initialiation:
	if self.tgbot, e = new(telegram.TelegramBot).ConfigureAndConnect(ctx); e != nil { return nil,e }

	return self,nil
}
func (self *application) LaunchEvLoops() error {
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

	return self.stopAndDestroy()
}
func (self *application) stopAndDestroy() error {
	return nil
}
func (self *application) SetLogger(log *zerolog.Logger) *application {
	self.log = log
	return self
}

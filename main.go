package main

import "os"
import "time"

import "mailru/rooster22/config"
import "mailru/rooster22/telegram"

import "github.com/rs/zerolog"


func main() {
	var e error
	var app *application = new(application)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	app.log = zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).With().Timestamp().Logger()
	app.log.Debug().Msg("Logger has been initialized!")

	// some initialization:
	if app.cfg, e = new(config.AppConfig).Parse(); e != nil {
		app.log.Error().Err(e).Msg("Application configuration error!")
		return
	}

	if _,e = new(telegram.TelegramBot).ConfigureAndConnect(app.cfg.Telegram.Token); e != nil {
		app.log.Error().Err(e).Msg("Telegram component error!")
		return
	}

	app.log.Debug().Msg("Application has been started and initialized!")
	app.log.Debug().Str(app.cfg.Mysql.Host, app.cfg.Telegram.Token).Msg("test")

	time.Sleep(60 * time.Second)
}

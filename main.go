package main

import "os"

import "mailru/rooster22/config"

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

	app.log.Debug().Msg("Application has been started and initialized!")
	app.log.Debug().Str(app.cfg.Mysql.Host, app.cfg.Telegram.Token).Msg("test")
}

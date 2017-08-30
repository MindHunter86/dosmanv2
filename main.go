package main

import "os"

import "github.com/rs/zerolog"


func main() {
	var e error
	var log zerolog.Logger

	// log iniitialization:
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log = zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).With().Timestamp().Logger()
	log.Debug().Msg("Logger has been initialized!")

	// application initialization:
	var app *application
	if app, e = new(application).SetLogger(&log).CreateAndConfigure(); e != nil {
		log.Error().Err(e).Msg("Error in the  application initialization!")
		os.Exit(1)
	}

	// launch application event loop:
	log.Debug().Msg("The Application has been initialized! Starting event loop...")
	if e = app.LaunchEvLoops(); e != nil {
		log.Error().Err(e).Msg("Error in application event loop launcher!")
		os.Exit(1)
	}

	// exit:
	log.Debug().Msg("The Application has been successfully stopped && destroyed!")
	os.Exit(0)
}

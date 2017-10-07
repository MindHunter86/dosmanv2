package main

import "os"
import "mailru/rooster22/system"
import "github.com/rs/zerolog"

func main() {
	var log zerolog.Logger

	// log iniitialization:
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log = zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).With().Timestamp().Logger()
	log.Debug().Msg("Logger has been initialized!")

	// system initialization:
	var e error
	var sys *system.System

	if sys,e = new(system.System).SetLogger(&log).Configure(); e != nil {
		log.Error().Err(e).Msg("Error in system initialization!")
		os.Exit(1)
	}

	// launch modules:
	if e = sys.Launch(); e != nil {
		log.Error().Err(e).Msg("Error in application launch handler!")
		os.Exit(1)
	} else { log.Debug().Msg("All modules has been launched!") }

	// launch system event loop:
	log.Debug().Msg("The Application has been initialized! Starting event loop...")
	if e = sys.LaunchEvLoops(); e != nil {
		log.Error().Err(e).Msg("Error in application event loop launcher!")
		os.Exit(1)
	}

	// exit:
	log.Debug().Msg("The Application has been successfully stopped && destroyed!")
	os.Exit(0)
}

package main

import "os"
import "mh00appserver/system"
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

	// system services start:
	log.Debug().Msg("Starting system bootstrapping...")
	if e = sys.Bootstrap(); e != nil {
		log.Error().Err(e).Msg("CRITICAL ERROR!")
		os.Exit(1)
	}

	// exit:
	log.Debug().Msg("The Application has been successfully stopped && destroyed!")
	os.Exit(0)
}

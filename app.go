package main

import "github.com/rs/zerolog"
import "mailru/rooster22/config"


type application struct {
	cfg *config.AppConfig
	log zerolog.Logger
}

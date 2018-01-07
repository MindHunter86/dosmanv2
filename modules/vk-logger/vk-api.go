package main

import "github.com/rs/zerolog"


type vkApi struct {
	log zerolog.Logger
}

func (m *vkApi) bootstrap() error { return nil }

package main

import (
	"mh00appserver/modules"
	"github.com/rs/zerolog"
)


type dbDriver struct { *modules.Modules }

func (m *dbDriver) configure(mods *modules.Modules) (dbDriver, error) {
	return &dbDriver{Modules: mods},nil
}

func (m *dbDriver) customerPut(customer) error {

	return nil
}

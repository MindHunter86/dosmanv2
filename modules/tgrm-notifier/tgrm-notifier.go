package main

import "reflect"

import "dosmanv2/modules"
import "github.com/rs/zerolog"

import "database/sql"
import _ "github.com/go-sql-driver/mysql"


type TgrmNotifier struct {
	db *sql.DB
	log *zerolog.Logger

	mod string
	mods *modules.Modules
}


// TgrmNotifier API:
func (m *TgrmNotifier) Construct(mods *modules.Modules, args ...interface{}) (modules.Module, error) {
	m.mods = mods
	m.mod = reflect.TypeOf(m).Elem().Name()
	m.log = m.mods.Logger.With().Str("plugin", m.mod.Logger()

	// initilize db connection:
	if e := m.dbInitialize(); e != nil { return nil,e }

	return m,nil
}

func (m *TgrmNotifier) Destruct() error { return nil }

func (m *TgrmNotifier) Bootstrap() error { return nil }


// TgrmNotifier internal API:
//

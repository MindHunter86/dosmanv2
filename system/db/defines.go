package db

import "database/sql"


type DBDriver interface {
	Construct(creds *DBCredentials) (DBDriver, error)
	Destruct() error

	GetRawDBSession() *sql.DB
}

type DBCredentials struct {
	Host, Username, Password, Database string
	MgrDirectory string
	MgrVersion uint
	Debug bool
}

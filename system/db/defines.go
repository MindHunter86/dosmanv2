package db

type DBDriver interface {
	Construct(creds *DBCredentials) (DBDriver, error)
	Destruct() error
}

type DBCredentials struct {
	Host, Username, Password, Database string
	MgrDirectory string
	MgrVersion uint
	Debug bool
}

package db

import (
	"os"
	"errors"
	"database/sql"
	"time"

	"github.com/rs/zerolog"
	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/go-sql-driver/mysql"

	"github.com/mattes/migrate"
	mysql_migrate "github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
)


var (
	errDBCredsIsNil = errors.New("Colud not exec construct for DBDriver: given credentials variable are not valid!")
)


// MySQLDriver definitions:
type MySQLDriver struct {
	log zerolog.Logger
	session *sql.DB
	credentials *DBCredentials
	migration *migrate.Migrate
}


// DBDriver API:
func (m *MySQLDriver) Construct(creds *DBCredentials) (DBDriver, error) {
	if creds == nil { return nil,errDBCredsIsNil }

	// log initialization for active debugging:
	if creds.Debug {
		m.log = zerolog.New(zerolog.ConsoleWriter{
			Out: os.Stderr }).With().Timestamp().Str("module", "DBDriver").Logger()
		m.log.Debug().Msg("DBDriver logger for active debug has been successfully initialized!")
	}

	// set default values for mysql credentials:
	switch {
		case creds.Database == "":
			creds.Database = "dosmanv2"
		case creds.Host == "":
			creds.Host = "localhost"
		case creds.Password == "":
			creds.Password = "1234"
		case creds.Username == "":
			creds.Username = "dosmanv2"
		case creds.MgrDirectory == "":
			creds.MgrDirectory = "migrations"
	}
	m.credentials = creds

	m.debugEcho(nil, "Check MySQL migrations ...")
	if mSession,err := sql.Open("mysql", m.configureConnection().FormatDSN()); err == nil {
		defer mSession.Close()

		if e := m.runMigrations(mSession); e != nil {
			m.debugEcho(e, "MySQL migrations were failed!")
			return nil,e
		}

		m.debugEcho(nil, "MySQL migrations are OK!")
	} else { return nil,err }

	return m,m.createConnection()
}

func (m *MySQLDriver) Destruct() error { return m.dropConnection() }


// MySQLDriver internal API:
func (m *MySQLDriver) debugEcho(e error, msg string) {
	if ! m.credentials.Debug { return }

	switch e {
		case nil: m.log.Debug().Msg(msg)
		default: m.log.Debug().Err(e).Msg(msg)
	}
}

func (m *MySQLDriver) createConnection() error {
	var e error
	if m.session,e = sql.Open("mysql", m.configureConnection().FormatDSN()); e != nil { return e }
	return m.session.Ping()
}

func (m *MySQLDriver) configureConnection() *mysql.Config {
	// https://github.com/go-sql-driver/mysql - mysql lib configuration

	location, e := time.LoadLocation("Europe/Moscow"); if e != nil {
		location = time.UTC
		m.debugEcho(e, "Could not get location for DBDriver configuration!")
	}

	return &mysql.Config{
		Net: "tcp",
		Addr: m.credentials.Host,
		User: m.credentials.Username,
		Passwd: m.credentials.Password,
		DBName: m.credentials.Database,
		Collation: "utf8_general_ci",
		MaxAllowedPacket: 0,
		TLSConfig: "false",
		Loc: location,

		Timeout: 10 * time.Second,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,

		AllowAllFiles: false,
		AllowCleartextPasswords: false,
		AllowNativePasswords: false,
		AllowOldPasswords: false,
		ClientFoundRows: false,
		ColumnsWithAlias: false,
		InterpolateParams: false,
		MultiStatements: true,
		ParseTime: true,
		Strict: m.credentials.Debug }
}

func (m *MySQLDriver) dropConnection() error { return m.session.Close() }

func (m *MySQLDriver) runMigrations(mSession *sql.DB) error {
	var e error

	mDriver, e := mysql_migrate.WithInstance(mSession, &mysql_migrate.Config{}); if e != nil { return e }
	if m.migration, e = migrate.NewWithDatabaseInstance("file://"+m.credentials.MgrDirectory, "mysql", mDriver); e != nil {
		return e
	}

	if e = m.migration.Up(); e != nil && e != migrate.ErrNoChange { return e }
	return nil
}

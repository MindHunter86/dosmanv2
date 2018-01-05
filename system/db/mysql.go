package db

import (
	"database/sql"
	"reflect"
	"time"

	"mh00appserver/modules"

	"github.com/rs/zerolog"
	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/go-sql-driver/mysql"
)


// Module structs:
type MySQLModule struct {
	dbSession *sql.DB
	migrations *sqlMigrate
	log zerolog.Logger
	topic *broker.Topic

	modName string
	mods *modules.Modules
}


// DBDriver API:
// func (m *MySQLModule) Construct(config)

// Module API:
func (self *MysqlModule) Configure(mods *modules.Modules, args ...interface{}) (modules.Module, error) {

	var e error

	self.mods = mods
	self.modName = reflect.TypeOf(self).Elem().Name()
	self.log = self.mods.Logger.With().Str("MODULE", self.modName).Logger()
	self.topic,e = self.mods.Broker.CreateTopic("mysql"); if e != nil { return nil,e }

	return self,self.openConnection()
}

func (self *MysqlModule) Bootstrap() error {
	var e error
	var brokerInbox chan *broker.Message
	var mysqlChecker *time.Ticker = time.NewTicker(time.Second)

	self.log.Debug().Msg("Check and Up mysql migrations ...")
	if sqlSession, err := sql.Open("mysql", self.configureConnetcion().FormatDSN()); err == nil { // connection for sql migrations
		if self.migrations, e = new(sqlMigrate).migrate(sqlSession, self.mods.Config.Mysql.Migrations.Dir, self.mods.Config.Mysql.Migrations.Version); e != nil {
			self.log.Error().Err(e).Msg("MySQL migrations error!")
			return e
		}
		self.log.Debug().Msg("MySQL migrations are OK!")
	} else { return err }

	brokerInbox = self.topic.Subscribe().GetInbox()

	self.log.Debug().Msg("Mysql has been bootstrapped!")
LOOP:
	for {
		select {
		case <-brokerInbox:
			self.log.Debug().Msg("brokerInbox has been triggered!")
		case <-self.mods.DonePipe:
			break LOOP
		case <-mysqlChecker.C:
			if _, e = self.dbSession.Exec("DO 1;"); e != nil { break LOOP }
		}
	}

	// stop timer and close mysql connection:
	mysqlChecker.Stop()
	if err := self.closeConnection(); err != nil {
		self.log.Error().Err(err).Msg("Could not successfully close mySQL connection!")
	}

	// return only for/select events:
	return e
}


// MysqlModule internal methods:
func (self *MysqlModule) openConnection() error {
	var e error
	if self.dbSession,e = sql.Open("mysql", self.configureConnetcion().FormatDSN()); e != nil { return e}
	return self.dbSession.Ping()
}

func (self *MysqlModule) closeConnection() error {
	return self.dbSession.Close()
}

func (self *MysqlModule) configureConnetcion() *mysql.Config {
	var cnf *mysql.Config = new(mysql.Config)

	// https://github.com/go-sql-driver/mysql - docs
	cnf.Net = "tcp"
	cnf.Addr = self.mods.Config.Mysql.Host
	cnf.User = self.mods.Config.Mysql.Username
	cnf.Passwd = self.mods.Config.Mysql.Password
	cnf.DBName = self.mods.Config.Mysql.Database
	cnf.Collation = "utf8_general_ci"
	cnf.MaxAllowedPacket = 0
	cnf.TLSConfig = "false"
	if tloc, e := time.LoadLocation("Europe/Moscow"); e != nil {	// "Europe%2FMoscow"
		self.log.Warn().Err(e).Msg("Could not get location in configuration files parsing!")
		//		self.log.W(log.LLEV_DBG, "Time location parsing error! | " + e.Error())
		cnf.Loc = time.UTC
	} else { cnf.Loc = tloc }

	cnf.Timeout = 10 * time.Second
	cnf.ReadTimeout = 5 * time.Second
	cnf.WriteTimeout = 10 * time.Second

	cnf.AllowAllFiles = false
	cnf.AllowCleartextPasswords = false
	cnf.AllowNativePasswords = false
	cnf.AllowOldPasswords = false
	cnf.ClientFoundRows = false
	cnf.ColumnsWithAlias = false
	cnf.InterpolateParams = false
	cnf.MultiStatements = true
	cnf.ParseTime = true
	cnf.Strict = true // XXX: Only for debug

	return cnf
}

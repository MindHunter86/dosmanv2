package mysql

import (
	"reflect"
	"time"
	"database/sql"

	"mailru/rooster22/modules"

	"github.com/rs/zerolog"
	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/go-sql-driver/mysql"
)


// Module structs:
type MysqlModule struct {
	dbSession *sql.DB
	log zerolog.Logger

	modName string
	mods *modules.Modules
}


// Module API:
func (self *MysqlModule) Configure(mods *modules.Modules, args ...interface{}) (modules.Module, error) {
	self.mods = mods
	self.modName = reflect.TypeOf(self).Elem().Name()

	// Set module name as prefix for logger:
	self.log = self.mods.Logger.With().Str("MODULE", self.modName).Logger()

	return self,self.openConnection()
}

func (self *MysqlModule) Bootstrap() error {
	var mysqlChecker *time.Ticker = time.NewTicker(time.Second)

LOOP:
	for {
		select {
		case <-self.mods.DonePipe:
			break LOOP
		case <-mysqlChecker.C:
			if _, e := self.dbSession.Exec("DO 1;"); e != nil {
				// XXX: garbage
				self.mods.ErrorPipe <-&modules.ModuleError{ ModName: self.modName, E: e }
				// self.log.Warn().Err(err).Msg("MysqlChecker error! Could not close the current connection!")

				self.log.Error().Err(e).Msg("MysqlChecker: mysql ping() has been failed!")
				// break LOOP - no break! Module interface will call close(), if we send error over ErrorPipe.
			}
		}
	}

	// stop timer:
	mysqlChecker.Stop()

	// close mysql connection:
	if e := self.closeConnection(); e != nil {
		self.log.Error().Err(e).Msg("Exception!")
	}

	return nil
}


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
	cnf.MultiStatements = false
	cnf.ParseTime = true
	cnf.Strict = true // XXX: Only for debug

	return cnf
}

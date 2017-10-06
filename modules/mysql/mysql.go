package mysql

import (
	"reflect"
	"time"
	"database/sql"

	"mailru/rooster22/modules"
	config "mailru/rooster22/system/config"

	"github.com/rs/zerolog"
	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/go-sql-driver/mysql"
)



type MysqlModule struct {
	dbSession *sql.DB

	log *zerolog.Logger
	cfg *config.SysConfig

	modName string
	mods *modules.Modules
}

func (self *MysqlModule) Configure(mods *modules.Modules, args ...interface{}) (modules.Module, error) {
	// Get and set module name (struct name):
	self.modName = reflect.TypeOf(self).Elem().Name()
	self.mods = mods

	// Get global logger and configuration:
	self.log = args[0].(*zerolog.Logger)
	self.cfg = args[1].(*config.SysConfig)

	go self.startCloseEventLoop()
	return self,self.openConnection()
}
func (self *MysqlModule) Start() error { return nil }
func (self *MysqlModule) Stop() error { return nil }
func (self *MysqlModule) Unconfigure() {}

func (self *MysqlModule) startCloseEventLoop() {
	<-self.mods.DonePipe
	self.mods.WaitGroup.Done()
	self.log.Debug().Msg("mysql - donePipe closed!")

	if e := self.closeConnection(); e != nil {
		self.log.Error().Err(e).Msg("Exception!")
	}
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
	cnf.Addr = self.cfg.Mysql.Host
	cnf.User = self.cfg.Mysql.Username
	cnf.Passwd = self.cfg.Mysql.Password
	cnf.DBName = self.cfg.Mysql.Database
	cnf.Collation = "utf8_general_ci"
	cnf.MaxAllowedPacket = 0
	cnf.TLSConfig = "false"
	if tloc, e := time.LoadLocation("Europe/Moscow"); e != nil {	// "Europe%2FMoscow"
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

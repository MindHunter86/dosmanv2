package mysql

import (
	"time"
	"database/sql"

	"mailru/rooster22/config"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/go-sql-driver/mysql"
)


type SqlConnector struct {
	dbConn *sql.DB

	appConfig *config.AppConfig
	appLogger *zerolog.Logger

	ctxPipeDone <-chan struct{}
}
func (self *SqlConnector) ConfigureAndConnect(ctx context.Context) error {
	self.ctxPipeDone = ctx.Done()
	self.appConfig = ctx.Value(config.CTX_APP_CONFIG).(*config.AppConfig)
	appLogger = ctx.Value(config.CTX_APP_LOGGER).(*zerolog.Logger)

	go self.startCloseEventLoop()
	return self.openConnection()
}

func (self *SqlConnector) startCloseEventLoop() {
	<-self.ctxPipeDone
	if e := self.closeConnection(); e != nil {
		self.appLogger.Error().Err(e).Msg("")
	}
}
func (self *SqlConnector) openConnection() error {
	if self.dbConn,e = sql.Open("mysql", self.configureConnetcion().FormatDSN()); e != nil { return e}
	return self.dbConn.Ping()
}
func (self *SqlConnector) closeConnection() error {
	return self.dbConn.Close()
}
func (self *SqlConnector) configureConnetcion() *mysql.Config {
	var cnf *mysql.Config = new(mysql.Config)

	// https://github.com/go-sql-driver/mysql - docs
	cnf.Net = "tcp4"
	cnf.Addr = self.appConfig.Mysql.Host
	cnf.User = self.appConfig.Mysql.Username
	cnf.Passwd = self.appConfig.Mysql.Password
	cnf.DBName = self.appConfig.Mysql.Database
	cnf.Collation = "utf8_general_ci"
	cnf.MaxAllowedPacket = 0
	cnf.TLSConfig = "false"
	if tloc, e := time.LoadLocation("Europe/Moscow"); e != nil {	// "Europe%2FMoscow"
		self.slogger.W(log.LLEV_DBG, "Time location parsing error! | " + e.Error())
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

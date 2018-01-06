package main

import "sync"
import "time"
import "reflect"
import "dosmanv2/modules"
import "dosmanv2/system/db"
import "github.com/rs/zerolog"


// XXX TEMPORARY CODE ZONE:
var maxWorkers int = 1
var workerJobTimeout time.Duration = 10 * time.Second
var workerJobTestPage string = "http://188.165.198.98:8089/"
var workerJobTestResponse []byte = []byte("Hello world\n")
var proxyCheckerOverdue uint = uint(300)
var proxyCheckerTimer time.Duration = 30 * time.Second
// 2DELETE END



// TODO:
// WORKER:
// - sync.waitgroup
// - worker.inbox as buffered chan
// DISPATCHER:
// - make chan proxy as send only chan




// ProxyChecker definitions:
var Plugin ProxyChecker

type ProxyChecker struct {
	db db.DBDriver
	log zerolog.Logger

	prxQueue chan proxy
	dispatcher *dispatcher

	proxyapi *proxyapi

	modName string
	mods *modules.Modules
	donePipe chan struct{}
}


// ProxyChecker API:
func (m *ProxyChecker) Construct(mods *modules.Modules, args ...interface{}) (modules.Module, error) {

	m.mods = mods
	m.donePipe = mods.DonePipe
	m.modName = reflect.TypeOf(m).Elem().Name()
	m.log = m.mods.Logger.With().Str("plugin", m.modName).Logger()

	var e error

	// initilize db connection:
	if e = m.dbInitialize(); e != nil { return nil,e }

	// initilize dispathcer with new proxy queue chan:
	m.prxQueue = make(chan proxy)
	m.dispatcher = &dispatcher{
		db: m.db,
		log: m.log,
		kernelQuit: m.donePipe,
		proxyQueue: m.prxQueue,
		pool: make(chan chan proxy, maxWorkers),
		workerQuit: make(chan struct{}, 1)}

	// initilize proxy api:
	m.proxyapi = &proxyapi{
		db: m.db.GetRawDBSession(),
		log: m.log,
		kernelQuit: m.donePipe,
		proxyCheckQueue: m.prxQueue}
	return m,nil
}

func (m *ProxyChecker) Bootstrap() error {
	var wg sync.WaitGroup

	m.log.Debug().Msg("Trying to bootstrap proxyapi...")
	go func(wg sync.WaitGroup) { wg.Add(1); m.proxyapi.bootstrap(); wg.Done() }(wg)

	m.log.Debug().Msg("Trying to bootstrap dispatcher...")
	go func(wg sync.WaitGroup) { wg.Add(1); m.dispatcher.bootstrap(m.proxyapi); wg.Done() }(wg)

	<-m.donePipe

	wg.Wait()
	m.log.Debug().Msg("Bootstrap func has been successfully completed!")
	return nil
}

func (m *ProxyChecker) Destruct() error {
	// close db connections:
	if e := m.db.Destruct(); e != nil { m.log.Error().Err(e).Msg("Could not successfully close DB connections!") }

	return nil
}


// ProxyChecker internal API:
func (m *ProxyChecker) dbInitialize() error {
	var e error

	if m.db,e = new(db.MySQLDriver).Construct(&db.DBCredentials{
		Host: m.mods.Config.Mysql.Host,
		Username: m.mods.Config.Mysql.Username,
		Password: m.mods.Config.Mysql.Password,
		Database: m.mods.Config.Mysql.Database,
		Debug: m.mods.Debug,
	}); e != nil { return e }

	return nil
}

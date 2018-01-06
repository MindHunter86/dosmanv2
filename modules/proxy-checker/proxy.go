package main

import "sync"
import "github.com/rs/zerolog"

import "time"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"


type proxy struct {
	addr string
	class, anon uint8
}
type proxyReport struct {
	proxy *proxy
	state bool
}

type proxies struct {
	wg sync.WaitGroup

	db *sql.DB
	log zerolog.Logger

	proxyCheckQueue chan proxy
}


func (m *proxies) construct(argDB *sql.DB, argLog zerolog.Logger, argQueue chan proxy) error {
	m.db = argDB
	m.log =argLog
	m.proxyCheckQueue = argQueue

	return nil
}

func (m *proxies) bootstrap(quit chan struct{}) {

	var isTheEnd, isRunning *bool
	var t1 *time.Timer = time.NewTimer(30 * time.Second)

	*isRunning = true // TODO : Use this variable!
	go m.readAllNonChecked(isTheEnd)

LOOP:
	for {
		select {
		case <-t1.C:
			go m.readAllOverdue(isTheEnd, proxyCheckerOverdue)
		case <-quit:
			break LOOP
		}
	}

	*isTheEnd = true
	m.wg.Wait()
}

func (m *proxies) readAllNonChecked(isTheEnd *bool) {
	m.wg.Add(1)
	defer m.wg.Done()

	m.log.Debug().Msg("ReadAllNonCheckd function is working now...")

	rows,e := m.db.Query("SELECT addr FROM proxies WHERE check = NULL"); if e != nil {
		m.log.Error().Err(e).Msg("Could not fetch records from database!")
		return
	}

	for rows.Next() {
		if *isTheEnd { break }

		var prx *proxy = new(proxy)
		if e := rows.Scan(&prx.addr); e != nil {
			m.log.Warn().Err(e).Msg("Could not scan the fetched row!")
			continue
		}

		m.proxyCheckQueue<- *prx
	}

	rows.Close()

	m.log.Debug().Msg("ReadAllNonCheckd function has been successfully completed!")
}

func (m *proxies) readAllOverdue(isTheEnd *bool, overdue uint) {
	m.wg.Add(1)
	defer m.wg.Done()

	m.log.Debug().Msg("ReadAllOverdue function is working now...")

	stmt,e := m.db.Prepare("SELECT p.addr FROM proxies p INNER JOIN checks c ON p.check = c.id WHERE TIMESTAMPDIFF(SECOND, c.checktime, NOW()) > ?")
	if e != nil { m.log.Error().Err(e).Msg("Could not prepare db statement for record fetching!"); return }
	defer stmt.Close()

	rows,e := stmt.Query(overdue); if e != nil { m.log.Error().Err(e).Msg("Could not fetch records from database!"); return }
	defer rows.Close()

	for rows.Next() {
		if *isTheEnd { break }

		var prx *proxy = new(proxy)
		if e := rows.Scan(&prx.addr); e != nil {
			m.log.Warn().Err(e).Msg("Could not scan the fetched row!")
			continue
		}

		m.proxyCheckQueue<- *prx
	}

	m.log.Debug().Msg("ReadAllOverdue function has been successfully completed!")
}

func (m *proxies) writeCheckerReport(report *proxyReport) error {
	stmt,e := m.db.Prepare("INSERT INTO checks (proxy, state) VALUES (?, ?)"); if e != nil { return e }
	res,e := stmt.Exec(report.proxy.addr, report.state); if e != nil { return e }
	rowId,e := res.LastInsertId(); if e != nil { return e }

	m.log.Debug().Int64("db_record_id", rowId).Msg("New checker report has been successfully writed in database!")
	return nil
}

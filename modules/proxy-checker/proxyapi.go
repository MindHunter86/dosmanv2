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

type proxyapi struct {
	wg sync.WaitGroup

	db *sql.DB
	log zerolog.Logger
	kernelQuit chan struct{}

	proxyCheckQueue chan proxy
}


func (m *proxyapi) bootstrap() {

	var isTheEnd *bool = new(bool)
	var t1 *time.Ticker = time.NewTicker(proxyCheckerTimer)

// var isRunning bool = true // TODO : Use this variable!
	go m.readAllNonChecked(isTheEnd)

LOOP:
	for {
		select {
		case <-t1.C: go m.readAllOverdue(isTheEnd, proxyCheckerOverdue)
		case <-m.kernelQuit: break LOOP
		}
	}

	*isTheEnd = true
	m.wg.Wait()
}

func (m *proxyapi) readAllNonChecked(isTheEnd *bool) {
	m.wg.Add(1)
	defer m.wg.Done()

	m.log.Debug().Msg("ReadAllNonCheckd function is working now...")

	rows,e := m.db.Query("SELECT `addr` FROM `proxies` WHERE `check` IS NULL"); if e != nil {
		m.log.Error().Err(e).Msg("Could not fetch records from database!")
		return
	}

	m.log.Debug().Bool("isTheEnd", *isTheEnd).Msg("Temporary debug")

	for rows.Next() {
		if *isTheEnd { break }

		var prx *proxy = new(proxy)
		if e := rows.Scan(&prx.addr); e != nil {
			m.log.Warn().Err(e).Msg("Could not scan the fetched row!")
			continue
		}

		m.log.Debug().Str("host", prx.addr).Msg("Send new job for checker...")
		m.proxyCheckQueue<- *prx
	}

	if e = rows.Err(); e != nil {
		m.log.Error().Err(e).Msg("Rows.Err() has been triggered!")}

	rows.Close()

	m.log.Debug().Msg("ReadAllNonCheckd function has been successfully completed!")
}

func (m *proxyapi) readAllOverdue(isTheEnd *bool, overdue uint) {
	m.wg.Add(1)
	defer m.wg.Done()

	m.log.Debug().Msg("ReadAllOverdue function is working now...")

	stmt,e := m.db.Prepare("SELECT p.addr FROM proxies p INNER JOIN `checks` c ON p.check = c.id WHERE TIMESTAMPDIFF(SECOND, c.checktime, NOW()) > ? AND p.check IS NOT NULL")
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

		m.log.Debug().Str("host", prx.addr).Msg("Send new job for checker...")
		m.proxyCheckQueue<- *prx
	}

	if e = rows.Err(); e != nil {
		m.log.Error().Err(e).Msg("Rows.Err() has been triggered!"); return }

	m.log.Debug().Msg("ReadAllOverdue function has been successfully completed!")
}

func (m *proxyapi) writeCheckerReport(report *proxyReport) error {
	m.log.Debug().Str("addr", report.proxy.addr).Msg("Started report writing!")

	stmt,e := m.db.Prepare("INSERT INTO checks(proxy, state) VALUES (?, ?)"); if e != nil { return e }
	res,e := stmt.Exec(report.proxy.addr, report.state); if e != nil { return e }
	rowId,e := res.LastInsertId(); if e != nil { return e }

	m.log.Debug().Str("addr", report.proxy.addr).Msg("insert into OK")

	stmt,e = m.db.Prepare("UPDATE `proxies` SET `check` = ? WHERE `addr` = ?"); if e != nil { return e }
	defer stmt.Close()

	res,e = stmt.Exec(rowId, report.proxy.addr); if e != nil { return e }
	if rwAff,e := res.RowsAffected(); e == nil && rwAff != 1 {
		m.log.Warn().Int64("rows_affected", rwAff).Msg("WriteChecherReport found some unstable transactions! (update proxies set check = id)")
	} else if e != nil { return e }

	m.log.Debug().Int64("db_record_id", rowId).Msg("New checker report has been successfully writed in database!")
	return nil
}

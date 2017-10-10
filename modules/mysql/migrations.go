package mysql

import "database/sql"
import "github.com/pressly/goose"
import _ "github.com/go-sql-driver/mysql"


type dbMigrations struct {
	db *sql.DB
	dir string
	migrs goose.Migrations
}


func (m *dbMigrations) configure(db *sql.DB, dir string) (*dbMigrations, error) {
	m.db = db
	m.dir = dir

	var e error
	if e = goose.SetDialect("mysql"); e != nil { return nil,e }
	if m.migrs, e = goose.CollectMigrations(m.dir, 0, goose.MaxVersion); e != nil {
		return nil,e
	} else { return m,nil }
}

func (m *dbMigrations) upMigrations() error {
	var row goose.MigrationRecord

	for _, mgr := range m.migrs {
		var query string = "select is_applied from goose_db_version where version_id=? order by tstamp desc limit 1"

		if e := m.db.QueryRow(query, mgr.Version).Scan(&row.IsApplied); e != nil && e == sql.ErrNoRows {
			if err := mgr.Up(m.db); err != nil { return err }
		} else if e != nil { return e }
	}

	return nil
}

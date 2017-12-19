package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "github.com/mattes/migrate"
import "github.com/mattes/migrate/database/mysql"
import _ "github.com/mattes/migrate/source/file"


type sqlMigrate struct { *migrate.Migrate }
func (m *sqlMigrate) migrate(db *sql.DB, dir string, ver uint) (*sqlMigrate, error) {

	var e error

	sqlDriver, _ := mysql.WithInstance(db, &mysql.Config{})
	if m.Migrate, e = migrate.NewWithDatabaseInstance("file://"+dir, "mysql", sqlDriver); e != nil {
		return nil,e
	}

	if e = m.Migrate.Migrate(ver); e != nil && e != migrate.ErrNoChange { return nil,e }
	return m,nil
}

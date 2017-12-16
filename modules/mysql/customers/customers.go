package customers


import "database/sql"
import _ "github.com/go-sql-driver/mysql"
//import mysql "github.com/go-sql-driver/mysql"

type Customers struct {
	db *sql.DB
	Pool CustomerPool
}

func (m *Customers) Configure(db *sql.DB) (*Customers, error) {
	if e := db.Ping(); e != nil {
		m.db = db
	} else { return nil,e }

	m.Pool = new(baseCustomerPool).configure(db)
	return m,nil
}

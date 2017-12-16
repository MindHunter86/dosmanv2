package customers

import "time"
import "sync"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
// import mysql "github.com/go-sql-driver/mysql"


type Customer struct {
	Phone, Firstname, Lastname string
	Chatroom int
	Registration *time.Time
}
type CustomerPool interface {
	GetCustomer(string) (Customer, error)
//	PutCustomer(string, *Customer) error
//	DropCustomer(string) error
}
type baseCustomerPool struct {
	db *sql.DB

	sync.RWMutex
	pool map[string]Customer // we used phone number as key for this pool
}


func (m *baseCustomerPool) configure(db *sql.DB) CustomerPool {
	m.db = db
	m.pool = make(map[string]Customer)
	return m
}
func (m *baseCustomerPool) GetCustomer(phone string) (Customer, error) {
	var customer *Customer = new(Customer)
	if e := m.db.QueryRow("SELECT * FROM customers WHERE phone=? LIMIT 1", phone).Scan(&customer); e != nil {
		if e == sql.ErrNoRows { return Customer{},nil }
		return Customer{},e
	}

	m.Lock()
	m.pool[phone] = *customer
	m.Unlock()

	return *customer,nil
}
//func (m *baseCustomerPool) PutCustomer(phone string, customer *Customer) error { return nil }
//func (m *baseCustomerPool) DropCustomer(phone string) error { return nil }
/*

func (m *TgrmMysqlApi) UpdateCustomer(phone, field string, value interface{}) error {
	var customer *TgrmCustomer = new(TgrmMysqlApi).GetCustomer(phone)
	mField,ok := reflect.TypeOf(customer).FieldByName(filed); if !ok { return errors.New("reflect error") }

	var mValue *reflect.Value := reflect.New(mField.Type)
	mValue.Set(value)

	if e := m.db.Query("UPDATE customers SET ?=? where phone=?", mField.Name, value, phone); e != nil { return e }


//		var query string = "select is_applied from goose_db_version where version_id=? order by tstamp desc limit 1"
//
//		if e := m.db.QueryRow(query, mgr.Version).Scan(&row.IsApplied); e != nil && e == sql.ErrNoRows {
//			if err := mgr.Up(m.db); err != nil { return err }
//		} else if e != nil { return e }


	return nil
}
*/

package mysql

import (
	"time"
	"errors"
	"reflect"
	"database/sql"

	"github.com/rs/zerolog"
	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/go-sql-driver/mysql"
)

type TgrmMysqlApi struct {
	db *sql.DB
}

type TgrmCustomer struct {
	Phone, Firstname, Lastname string
	Chatroom int
	Registration *time.Time
}

// TgrmMysqlApi API:
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
func (m *TgrmMysqlApi) GetCustomer(phone string) (*TgrmCustomer, error) {}
func (m *TgrmMysqlApi) CreateCustomer() error {}


// TgrmMysqlApi initial methods:
func (m *tgrmMysqlApi) configure() {}


func (m *tgrmMysqlApi) getReflectType(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.S
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return nil
}

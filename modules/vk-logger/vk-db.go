package main

import "time"
import "errors"
import "net/http"
import "dosmanv2/system/config"
import "github.com/boltdb/bolt"
import "github.com/rs/zerolog"


type vkDB struct{
	db *bolt.DB
	log *zerolog.Logger
}


// vkDB internal API:
func (m *vkDB) construct(log *zerolog.Logger, config *config.SysConfig) (*vkDB, error) {
	m.log = log

	var e error
	if m.db,e = bolt.Open(config.Vklogger.Cookie_storage.Path, 0600, &bolt.Options{
		Timeout: 1 * time.Second}); e != nil { return nil,e }

	dbTx,e := m.db.Begin(true); if e != nil { return nil,e }
	defer dbTx.Rollback()

	if dbTx.Bucket([]byte("vk_session")) == nil {
		m.log.Warn().Msg("VkDB warning! BoltDB schema is not valid! Trying to initilize new schema...")
		if _,e = dbTx.CreateBucket([]byte("vk_session")); e != nil { return nil,e }
		if e = dbTx.Commit(); e != nil { return nil,e }
	}

	return m,e
}

func (m *vkDB) destruct() error { return m.db.Close() }

func (m *vkDB) updateCookies(cookies []*http.Cookie) error {
	tx,e := m.db.Begin(true); if e != nil { return e }
	defer tx.Rollback()

	var vkSessBucket = tx.Bucket([]byte("vk_session"))
	for _,v := range cookies {
		m.log.Debug().Msg("UpdateCookies - starting DB transaction...")

		if e = vkSessBucket.Put([]byte(v.Name), []byte(v.Value)); e != nil {
			m.log.Warn().Err(e).Msg("UpdateCookies could not write cookie in BoltDB Bucket!")
			continue
		}
	}

	if e = tx.Commit(); e != nil { return e }
	return nil
}
func (m *vkDB) getCookies() ([]*http.Cookie, error) {
	tx,e := m.db.Begin(false); if e != nil { return nil,e }
	defer tx.Rollback()

	var httpCookies = make([]*http.Cookie, 0)
	var vkSessBucket = tx.Bucket([]byte("vk_session"))

	if vkSessBucket.Stats().KeyN == 0 { return nil,errors.New("BoltDB bucket vk_session is empty!") }
	m.log.Debug().Int("keys_count", vkSessBucket.Stats().KeyN).Msg("BoltDB keys count from vk_session bucket.")

	if e = vkSessBucket.ForEach(func(k,v []byte) error {
		httpCookies = append(httpCookies, &http.Cookie{
			Name: string(k),
			Value: string(v),
			Secure: true,
			HttpOnly: true})

		return nil
	}); e != nil { return nil,e }

	return httpCookies,nil
}

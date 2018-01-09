package main

import (
	"sync"
	"context"
	"time"
	"bytes"
	"errors"
	"reflect"

	"io"
	"io/ioutil"

	"net/url"
	"net/http"
	"net/http/cookiejar"

	"dosmanv2/modules"
	"dosmanv2/system/config"

	"golang.org/x/net/html"
	"github.com/rs/zerolog"
)


// Plugin variables:
var Plugin VKLogger

type VKLogger struct {
	log zerolog.Logger

	vkHttpClient *http.Client
	vkStorage *vkDB
	vkApi *vkApi

	modName string
	mods *modules.Modules
}


// VKLogger plugin API:
func (m *VKLogger) Construct(mods *modules.Modules, args ...interface{}) (modules.Module, error) {
	m.mods = mods
	m.modName = reflect.TypeOf(m).Elem().Name()
	m.log = m.mods.Logger.With().Str("plugin", m.modName).Logger()

	var e error
	if m.vkStorage,e = new(vkDB).construct(&m.log, m.mods.Config); e != nil { return nil,e }

	m.vkHttpClient = new(http.Client)

	if e = m.checkSavedSession(); e != nil {
		m.log.Warn().Err(e).Msg("Could not get save cookies from BoltDB db! Trying to create new cookieJar...")
		if err := m.vkAuthenticate(m.mods.Config); e != nil { return nil,err }
	}

	m.vkApi = new(vkApi).construct(&m.log, m.mods.Config, m)

	if _,e = m.vkGetWallPostWiget("-35005_29999"); e != nil { return nil,e } // call this method only for debugging
	return m,nil
}
func (m *VKLogger) Bootstrap() error {
	var wg sync.WaitGroup

	go m.vkApi.bootstrap(&wg)
	defer wg.Wait()

	m.log.Debug().Msg("Starting event loop...")
LOOP:
	for {
		select {
		case <-m.mods.DonePipe:
			m.log.Debug().Msg("Caught donePipe close!")
			if e := m.vkApi.httpSrv.Shutdown(context.Background()); e != nil {
				m.log.Error().Err(e).Msg("Could not initilize Shutdown method for http server!")
			}
			m.log.Debug().Msg("Caught donePipe close2!")
			break LOOP
		}
	}

	return nil
}
func (m *VKLogger) Destruct() error {
	m.log.Debug().Msg("VKLogger destruct method has been called!") // XXX: tmp debug
	return m.vkStorage.destruct()
}


// VKLogger plugin internal API:
func (m *VKLogger) checkSavedSession() error {
	vkURL,e := url.Parse("https://vk.com/"); if e != nil { return e }
	if m.vkHttpClient.Jar,e = cookiejar.New(nil); e != nil { return e }

	cookies,e := m.vkStorage.getCookies(); if e != nil { return e }
	m.vkHttpClient.Jar.SetCookies(vkURL, cookies)

	if m.mods.Debug {
		for _,v := range cookies { m.log.Debug().Str("loaded_cookie", v.String()).Msg("Found saved session cookie in BoltDB!") }}

	m.log.Info().Msg("CheckSavedSession - saved cookies have been successfully loaded and installed!")
	return nil
}
func (m *VKLogger) vkGetWallPostWiget(wallPostId string) (string,error) {
	var postBuf = new(bytes.Buffer)
	postBuf.WriteString("act=a_get_post_hash&al=1&post="+wallPostId)

	rsp,e := m.vkHttpClient.Post("https://vk.com/dev.php", "application/x-www-form-urlencoded", postBuf); if e != nil { return "",e }
	if rsp.StatusCode != 200 { m.log.Warn().Int("response_code", rsp.StatusCode).Msg("Method vkGetWallPostWiget found unstable response from VK! (url: https://vk.com/dev.php)") }
	defer rsp.Body.Close()

	rspBody,e := ioutil.ReadAll(rsp.Body); if e != nil { return "",nil }
	m.log.Info().Str("response_body", string(rspBody)).Msg("Method vkGetWallPostWiget get new response!")

	var rspSplit = bytes.Split(rspBody, []byte("<!>"))
	if ! bytes.Equal(rspSplit[4], []byte("0")) {
		return "",errors.New("Method vkGetWallPostWiget has unusual response from VK (https://vk.com/dev.php). We can not say more, because we don't know the VK proto specification.")
	}

	if m.mods.Debug {
		m.log.Info().Bytes("VALUE", rspSplit[5]).Msg("Found new post hash!")}
	return string(rspSplit[5]),nil
}

func (m *VKLogger) vkAuthenticate(config *config.SysConfig) error {
	var e error
	if m.vkHttpClient.Jar,e = cookiejar.New(nil); e != nil { return e }

	m.log.Debug().Msg("Tying to get VK main page...")
	rsp,e := m.vkHttpClient.Get("https://vk.com/"); if e != nil { return e }
	defer rsp.Body.Close()

	m.log.Debug().Msg("Trying to get hidden values from VK main page...")
	authTarget, authData := m.getFormHiddenValues(rsp.Body)
	if len(authTarget) == 0 { return errors.New("Authentication target is empty! Form parsing has been failed!") }

	authData.Add("email", config.Vklogger.Login)
	authData.Add("pass", config.Vklogger.Password)

	m.log.Debug().Msg("Trying to send POST request for vkAuthentication...")
	rsp,e = m.vkHttpClient.PostForm(authTarget, authData); if e != nil { return e }

	m.log.Debug().Msg("POST request has been sended! Trying to read response body...")
	rspBody,e := ioutil.ReadAll(rsp.Body); if e != nil { return e }
	if m.mods.Debug { m.log.Info().Msg(string(rspBody)) }

	authURL,e := url.Parse("https://vk.com/"); if e != nil { return e }
	if e = m.vkStorage.updateCookies(m.vkHttpClient.Jar.Cookies(authURL)); e != nil {
		m.log.Error().Err(e).Msg("Could not save given session cookie for logging!")
	}

	if m.mods.Debug {
		for _,v := range m.vkHttpClient.Jar.Cookies(authURL) {
			m.log.Info().Str("VALUE", v.String()).Msg("Found new cookie!")
		}}
	return nil
}

func (m *VKLogger) getFormHiddenValues(rspBody io.ReadCloser) (string, url.Values) {
	m.log.Debug().Msg("Tokenizer initialization...")
	tokenizer := html.NewTokenizer(rspBody)

	var formTarget string
	var formHiddenData = url.Values{}
	var parseTrack time.Time; if m.mods.Debug { parseTrack = time.Now()	}

	m.log.Debug().Msg("GetFormHiddenData parsing has been started!")
LOOP:
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			break LOOP

		case html.StartTagToken:
			switch token := tokenizer.Token(); token.Data {
			case "form":
				for _,attr := range token.Attr {
					if attr.Key == "action" { formTarget = attr.Val }
				}
			default: continue
			}

		case html.SelfClosingTagToken:
			switch token := tokenizer.Token(); token.Data {
			case "input":
				switch token.Attr[1].Val {
				case "act": formHiddenData.Add("act", token.Attr[2].Val)
				case "role": formHiddenData.Add("role", token.Attr[2].Val)
				case "_origin": formHiddenData.Add("_origin", token.Attr[2].Val)
				case "ip_h": formHiddenData.Add("ip_h", token.Attr[2].Val)
				case "lg_h": formHiddenData.Add("lg_h", token.Attr[2].Val)
				default: continue
				}

			default: continue
			}
		}
	}

	if m.mods.Debug {
		m.log.Debug().Dur("elapsed", time.Since(parseTrack)).Msg("GetFormHiddenData parser performance report.")}

	m.log.Debug().Msg("GetFormHiddenData parsing has been stopped!")
	return formTarget, formHiddenData
}

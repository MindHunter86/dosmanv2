package main

import "reflect"
import "github.com/rs/zerolog"
import "dosmanv2/modules"
import (
	"errors"
	"net/url"
	"net/http"
	"io/ioutil"
	"golang.org/x/net/html"
	"net/http/cookiejar"
	"io"
)

// Global variables for BitBucket export:
// n\a


// Plugin variables:
var Plugin VKLogger

type VKLogger struct {
  log zerolog.Logger

	vkHttpClient *http.Client

	modName string
	mods *modules.Modules
}


// VKLogger plugin API:
func (m *VKLogger) Construct(mods *modules.Modules, args ...interface{}) (modules.Module, error) {
	m.mods = mods
	m.modName = reflect.TypeOf(m).Elem().Name()
	m.log = m.mods.Logger.With().Str("plugin", m.modName).Logger()

	m.log.Debug().Msg("String Construct method for vk-logger plugin...")

	if e := m.vkAuthenticate("piratickgoo@gmail.com", "Ovij0jmNCSDWcCpl"); e != nil { return nil,e }

	m.log.Debug().Msg("Construct method stop for vk-logger plugin...")
	return m,nil
}
func (m *VKLogger) Bootstrap() error { return nil }
func (m *VKLogger) Destruct() error { return nil }


// VKLogger plugin internal API:
func (m *VKLogger) vkGetWallPostWiget(wallPostId string) error {

	

	return nil
}
func (m *VKLogger) vkAuthenticate(email, password string) error {
	var e error
	var httpClient *http.Client = new(http.Client)

	if httpClient.Jar,e = cookiejar.New(nil); e != nil { return e }

	m.log.Debug().Msg("Tying to get VK main page...")
	rsp,e := httpClient.Get("https://vk.com/"); if e != nil { return e }
	defer rsp.Body.Close()

	m.log.Debug().Msg("Trying to get hidden values from VK main page...")
	authTarget, authData := m.getFormHiddenValues(rsp.Body)
	if len(authTarget) == 0 { return errors.New("Authentication target is empty! Form parsing has been failed!") }

	authData.Add("email", email)
	authData.Add("pass", password)

	m.log.Debug().Msg("Trying to send POST request for vkAuthentication...")
	rsp,e = httpClient.PostForm(authTarget, authData); if e != nil { return e }

	m.log.Debug().Msg("POST request has been sended! Trying to read response body...")
	rspBody,e := ioutil.ReadAll(rsp.Body); if e != nil { return e }
	m.log.Info().Msg(string(rspBody))

	m.log.Debug().Msg("Final parsing...")
	authURL,e := url.Parse("https://vk.com/"); if e != nil { return e }
	for _,v := range httpClient.Jar.Cookies(authURL) {
		m.log.Info().Str("VALUE", v.String()).Msg("Found new cookie!")
	}

	return nil
}
func (m *VKLogger) getFormHiddenValues(rspBody io.ReadCloser) (string, url.Values) {
	m.log.Debug().Msg("Tokenizer initialization...")
	tokenizer := html.NewTokenizer(rspBody)

	var formTarget string
	var formHiddenData = url.Values{}

	m.log.Debug().Msg("Magic started")
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

	m.log.Debug().Msg("Magic stopped!")
	return formTarget, formHiddenData
}

/*
func (m *worker) doJob(api *proxyapi, prx *proxy) error {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL( &url.URL{ Host: prx.addr } )},
		Timeout: workerJobTimeout}

	rsp,e := httpClient.Get(workerJobTestPage); if e != nil { return e }
	defer rsp.Body.Close()

	rspBody,e := ioutil.ReadAll(rsp.Body); if e != nil { return e }
	log.Println(string(rspBody))
	return api.writeCheckerReport(&proxyReport{
		proxy: prx,
		state: bytes.Equal(rspBody, workerJobTestResponse)})
}*/

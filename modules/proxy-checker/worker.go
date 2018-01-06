package main

import "log"
import "bytes"
import "net/http"
import "net/url"
import "io/ioutil"

// TODO:
// - add User-Agent randomizer;


type worker struct {
	pool chan chan proxy
	inbox chan proxy
	quit chan struct{}
}

func (m *worker) construct(pool chan chan proxy, quit chan struct{}) *worker {
	m.pool = pool
	m.quit = quit
	m.inbox = make(chan proxy, 1) // XXX: track for it !!
	return m
}

func (m *worker) spawn(argProxyApi *proxyapi) {
	var e error

LOOP:
	for {
		m.pool <- m.inbox

		select {
		case prx := <-m.inbox:
			if e = m.doJob(argProxyApi, &prx); e != nil { log.Println(e) } // TODO : Report for error
		case <-m.quit:
			break LOOP
		}
	}
}

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
}

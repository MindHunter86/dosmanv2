package main

import "sync"


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

func (m *worker) spawn(wg *sync.WaitGroup) {
	wg.Add(1)

LOOP:
	for {
		m.pool <- m.inbox

		select {
		case _ = <-m.inbox:
			// TODO: CHECK PROXY
		case <-m.quit:
			break LOOP
		}
	}

	wg.Done()
}

//func (m *worker) proxyCheck(prx proxy) {
//	//
//}
//
//
//
//
//
//
//
////
//var (
//	errProxyIsEmpty = errors.New("Given proxy host is empty!")
//	errProxyAbnormalResult = errors.New("Cound not fetch test url!")
//)
//
//type ProxyDrvier struct {
//	log *zerolog.Logger
//}
//
//
//func (m *ProxyDrvier) Construct() (*ProxyDrvier, error) {
//	return m,m.proxyCheck("149.202.180.55:31288")
//}
//
//
//func (m *ProxyDrvier) proxyCheck(host string) error {
//	if len(host) == 0 { return errProxyIsEmpty }
//
//	var timeout = time.Duration(5 * time.Second)
//
//	var httpClient *http.Client = &http.Client{
//		Transport: &http.Transport{
//			Proxy: http.ProxyURL(&url.URL{
//				Host: host,
//			}),
//		},
//		Timeout: timeout}
//
//	resp, e := httpClient.Get("http://188.165.198.98:8089/"); if e != nil { return e }
//	defer resp.Body.Close()
//
//	body, _ := ioutil.ReadAll(resp.Body)
//	if ! bytes.Equal(body, []byte("Hello world\n")) {
//		log.Println(string(body))
//		log.Println(body)
//		return errProxyAbnormalResult
//	}
//
//	// TODO: Write proxy in mysql
//	return nil
//}

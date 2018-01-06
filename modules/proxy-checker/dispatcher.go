package main

import "sync"
import "dosmanv2/system/db"
import "github.com/rs/zerolog"


type dispatcher struct {
	db db.DBDriver
	log zerolog.Logger

	pool chan chan proxy
	proxyQueue chan proxy
	kernelQuit chan struct{}
	workerQuit chan struct{}
}


func (m *dispatcher) bootstrap(argProxyApi *proxyapi) {
	var wg sync.WaitGroup
	wg.Add(maxWorkers + 1)

	for i := 0; i < maxWorkers; i++ {
		go func(wg *sync.WaitGroup) { new(worker).construct(m.pool, m.workerQuit).spawn(argProxyApi); wg.Done() }(&wg)
	}

	go func(wg *sync.WaitGroup) { m.dispatch(); wg.Done() }(&wg)
	wg.Wait()
}

func (m *dispatcher) dispatch() {
	m.log.Debug().Msg("Dispatcher has been started!")
LOOP:
	for {
		select {
		case i := <-m.proxyQueue:
			m.log.Info().Str("proxy", i.addr).Msg("Dispatcher has a new job!")
			go func(prx proxy){ // XXX: Goroutine here??? Optimize!
				jobQueue := <-m.pool
				jobQueue <- prx
			}(i)
		case <-m.kernelQuit: break LOOP
		}
	}

	close(m.workerQuit)
}

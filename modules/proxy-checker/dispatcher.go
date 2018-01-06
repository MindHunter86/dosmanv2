package main

import "sync"
import "dosmanv2/system/db"
import "github.com/rs/zerolog"


type dispatcher struct {
	wg sync.WaitGroup

	db db.DBDriver
	log zerolog.Logger

	pool chan chan proxy
	proxyQueue chan proxy
	kernelQuit chan struct{}
	workerQuite chan struct{}
}


func (m *dispatcher) bootstrap() {
	for i := 0; i < maxWorkers; i++ {
		go new(worker).construct(m.pool, m.workerQuite).spawn(&m.wg)
	}

	m.dispatch()
}

func (m *dispatcher) dispatch() {
LOOP:
	for {
		select {
		case i := <-m.proxyQueue:
			go func(prx proxy){ // XXX: Goroutine here??? Optimize!
				jobQueue := <-m.pool
				jobQueue <- prx
			}(i)
		case <-m.kernelQuit: break LOOP
		}
	}

	close(m.workerQuite)
	m.wg.Wait()
}

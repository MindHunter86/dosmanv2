package main

import "sync"


type dispatcher struct {
	pool chan chan proxy
	proxyQueue chan proxy
	kernelQuit chan struct{}
	workerQuite chan struct{}

	wg sync.WaitGroup
}

func (m *dispatcher) construct(sigpipe chan struct{}, prxqueue chan proxy) *dispatcher {
	m.kernelQuit = sigpipe
	m.proxyQueue = prxqueue

	m.workerQuite = make(chan struct{}, 1)
	m.pool = make(chan chan proxy, maxWorkers)
	return m
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

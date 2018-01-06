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
		case i := <-m.inbox:
		case <-m.quit:
			break LOOP
		}
	}

	wg.Done()
}

package main


import "log" // XXX: temporary!!
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
	log.Println("Worker ### has been spawned!")

LOOP:
	for {
		m.pool <- m.inbox

		select {
		case i := <-m.inbox:
			log.Println("=== new job:")
			log.Println(i.addr)
			log.Println(i.anon)
			log.Println(i.class)
			log.Println(i.created.String())
			log.Println("== job end")
		case <-m.quit:
			break LOOP
		}
	}

	log.Println("Worker ### has been killed!")
	wg.Done()
}

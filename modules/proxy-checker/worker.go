package main


import "log" // XXX: temporary!!


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

func (m *worker) spawn() {
	log.Println("Worker ### has been spawned!")

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
			log.Println("Worker ### has been killed!")
			return
		}
	}
}

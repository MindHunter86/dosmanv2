package main


type dispatcher struct {
	pool chan chan proxy
	quitPipe chan struct{}
}

func (m *dispatcher) construct() *dispatcher {
	m.pool = make(chan chan proxy, maxWorkers)
	m.quitPipe = make(chan struct{}, 1)
	return m
}

func (m *dispatcher) bootstrap(queue chan proxy) {
	for i := 0; i < maxWorkers; i++ {
		go new(worker).construct(m.pool, m.quitPipe).spawn()
	}

	go m.dispatch(queue)
}

func (m *dispatcher) dispatch(queue chan proxy) {
	for {
		select {
		case i := <-queue:
			go func(prx proxy){ // XXX: Goroutine here??? Optimize!
				jobQueue := <-m.pool
				jobQueue <- prx
			}(i)
		}
	}
}

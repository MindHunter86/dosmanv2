package main

import "log"
import "time"
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

	go m.dispatch()

	time.Sleep(5*time.Second)

	var proxyTestRecord *proxy = &proxy{
		addr: "127.0.0.1:3128",
		class: uint8(1),
		anon: uint8(1),
		created: time.Now()}

	m.proxyQueue<- *proxyTestRecord

	m.wg.Wait()
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
	log.Println("quitPipe has been closed! sync.Wait...")
//	m.wg.Wait()
}

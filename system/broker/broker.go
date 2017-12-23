package broker

import (
	"sync"
	"errors"

	config "mh00appserver/system/config"
)


type Broker struct {
	wg sync.WaitGroup
	cnf *config.SysConfig

	sync.RWMutex
	topics map[string]*Topic
}

type Message func(interface{}) error
type Subscription struct { ch chan<- *Message }
type Topic struct { sub *Subscription }


// global error definitions:
var (
	ErrTopicIsDefined = errors.New("Could not create new topic; It's has already defined!")
	ErrTopicIsNotDefined = errors.New("Could not find topic!")
	ErrTopicIsClosed = errors.New("Could not push message into channel: Topic has been closed!")
)


// Topic external methods:
func (m *Topic) Subscribe() *Subscription {
	return m.sub
}

// Topic internal methods:
func (m *Topic) create(buf int) (*Topic, error) {

	m.sub = new(Subscription)
	m.sub.ch = make(chan<- *Message)

	return m,nil
}


// Subscription external methods:
func (m *Subscription) Publish(mess *Message) error {
	// XXX:	if _, ok := <-m.ch; !ok { return ErrTopicIsClosed } // read0only channel

	m.ch<- mess; return nil
}


// Broker external methods:
func (m *Broker) CreateTopic(name string) error {
	m.RLock(); _, ok := m.topics[name]; m.RUnlock()
	if ok { return ErrTopicIsDefined }

	m.Lock();
	defer m.Unlock()

	var e error
	if m.topics[name], e = new(Topic).create(m.cnf.Base.Broker.Buffer); e != nil { return e }

	return nil
}
func (m *Broker) DeleteTopic(name string) error {
	m.RLock(); _, ok := m.topics[name]; m.RUnlock()
	if ! ok { return ErrTopicIsNotDefined }

	m.Lock(); delete(m.topics, name); m.Unlock()
	return nil
}

// Broker internal methods:
func (m *Broker) Configure(cnf *config.SysConfig) (*Broker, error) {

	// define internal variables:
	m.cnf = cnf
	m.topics = make(map[string]*Topic)

	return m,nil
}

// TODO: DROP ALL TOPICS; CLOSE ALL CHANNELS !!!
func (m *Broker) Shutdown() {}

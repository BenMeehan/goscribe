package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// Arilator (aka Broker)
// Topics : map of all the currently active topics with atleast 1 subscriber
// Command : channel to recive the commands from clients
type Arilator struct {
	topics   map[string]*Topic
	commands chan Command
}

// Inititalize a new broker
func newArilator() *Arilator {
	return &Arilator{
		topics:   make(map[string]*Topic),
		commands: make(chan Command),
	}
}

// Keep listening on the command channel for new commands
func (a *Arilator) start() {
	for cmd := range a.commands {
		switch cmd.id {
		case PUB:
			a.publish(cmd.emptor, cmd.args)
		case SUB:
			a.subscribe(cmd.emptor, cmd.args)
		case UNSUB:
			a.unsubscribe(cmd.emptor, cmd.args)
		case LIST:
			a.list(cmd.emptor)
		case QUIT:
			a.disconnect(cmd.emptor)
		}
	}
}

// helper method to create a new client
func (a *Arilator) newClient(conn net.Conn) {
	log.Println("new emptor connected :", conn.RemoteAddr())
	e := &Emptor{
		conn:     conn,
		topics:   make(map[string]*Topic),
		commands: a.commands, // reference to the brokers commands channel
	}
	e.readMessage()
}

// ---------- BUSSINESS LOGIC METHODS -----------
func (a *Arilator) publish(e *Emptor, args []string) {
	if len(args) < 2 {
		e.Err(fmt.Errorf("not enough arguments"))
		return
	}
	topic := args[1]
	t, ok := a.topics[topic]
	if !ok {
		e.Warn("no subscribers for topic")
		return
	}
	if len(args) < 3 {
		e.Err(fmt.Errorf("not enough arguments"))
		return
	}
	for key := range t.subscribers {
		if key == e {
			continue
		}
		go key.Msg(topic + " " + args[2])
	}
	e.Msg("published message " + args[2] + " to topic " + topic)
}

func (a *Arilator) subscribe(e *Emptor, args []string) {
	if len(args) < 2 {
		e.Err(fmt.Errorf("not enough arguments"))
		return
	}
	topic := args[1]
	t, ok := a.topics[topic]
	if !ok {
		t = &Topic{
			name:        topic,
			subscribers: make(map[*Emptor]bool),
		}
		a.topics[topic] = t
	}
	t.subscribers[e] = true
	if _, ok = e.topics[t.name]; !ok {
		e.topics[t.name] = t
	}
	e.Msg("subscribed to topic " + topic)
}

func (a *Arilator) unsubscribe(e *Emptor, args []string) {
	if len(args) < 2 {
		e.Err(fmt.Errorf("not enough arguments"))
		return
	}
	topic := args[1]
	t, ok := a.topics[topic]
	if !ok {
		e.Err(fmt.Errorf("topic does not exist"))
		return
	}
	_, ok = t.subscribers[e]
	if !ok {
		e.Err(fmt.Errorf("not subscribed to that topic"))
		return
	}
	delete(t.subscribers, e)
	if len(t.subscribers) == 0 {
		delete(a.topics, topic)
	}
	e.Msg("unsubscribed from topic " + topic)
}

func (a *Arilator) list(e *Emptor) {
	var topics []string
	for key := range a.topics {
		topics = append(topics, key)
	}
	e.Msg(strings.Join(topics, ", "))
}

func (a *Arilator) disconnect(e *Emptor) {
	for key := range e.topics {
		t := a.topics[key]
		delete(t.subscribers, e)
		if len(t.subscribers) == 0 {
			delete(a.topics, key)
		}
	}
	e.Msg("bye")
	e.conn.Close()
}

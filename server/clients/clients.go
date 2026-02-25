package clients

import (
	"fmt"
	"net"
	"sync"
)

func New(conn net.Conn) *Client {
	c := makeClient(conn)

	Clients.mu.Lock()
	Clients.clients[c.RemoteAddr.String()] = &c
	Clients.mu.Unlock()

	Clients.events <- Event{
		Status:     Connected,
		RemoteAddr: c.RemoteAddr,
	}

	fmt.Println("debug: new client connected")
	return &c
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		Clients.mu.Lock()

		delete(Clients.clients, c.RemoteAddr.String())

		Clients.events <- Event{
			Status:     Disconnected,
			RemoteAddr: c.RemoteAddr,
		}

		Clients.mu.Unlock()

		close(c.done)
		fmt.Println("debug: client destroyed")
	})
}

func makeClient(conn net.Conn) Client {
	done := make(chan struct{})

	return Client{
		RemoteAddr: conn.RemoteAddr(),
		Conn:       conn,
		Done:       done,
		done:       done,
	}
}

var events chan Event = make(chan Event, 64)
var Events <-chan Event = events

var Clients = &clients{
	clients: make(map[string]*Client),
	events:  events,
	Events:  Events,
}

// types
type Status int

const (
	Connected Status = iota
	Disconnected
)

type Event struct {
	Status     Status
	RemoteAddr net.Addr
}

type clients struct {
	mu      sync.RWMutex
	clients map[string]*Client
	Events  <-chan Event
	events  chan Event
}

type Client struct {
	Conn       net.Conn
	Done       <-chan struct{}
	RemoteAddr net.Addr
	done       chan struct{}
	closeOnce  sync.Once
}

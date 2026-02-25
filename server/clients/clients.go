package clients

import (
	"errors"
	"net"
	"sync"
)

func Connect(conn net.Conn) Client {
	c := makeClient(conn)

	state.mu.Lock()
	state.clients[c.RemoteAddr.String()] = &c
	state.mu.Unlock()

	state.events <- statusEvent{
		status:     Connected,
		remoteAddr: c.RemoteAddr,
	}

	return c
}

func Disconnect(addr net.Addr) error {
	c, err := GetClient(addr)
	if err != nil {
		return err
	}

	state.mu.Lock()

	delete(state.clients, c.RemoteAddr.String())

	state.events <- statusEvent{status: Disconnected, remoteAddr: c.RemoteAddr}
	state.mu.Unlock()

	close(c.done)

	return nil
}

func GetClient(addr net.Addr) (*Client, error) {
	state.mu.RLock()
	c := state.clients[addr.String()]
	state.mu.RUnlock()

	if c == nil {
		return nil, errors.New("Could not find client")
	}
	return c, nil

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

var state = &clients{
	clients: make(map[string]*Client),
	events:  make(chan statusEvent),
}

// types
type Status int

const (
	Connected Status = iota
	Disconnected
)

type statusEvent struct {
	status     Status
	remoteAddr net.Addr
}

type Client struct {
	Conn       net.Conn
	Done       <-chan struct{}
	RemoteAddr net.Addr
	done       chan struct{}
}

type clients struct {
	mu      sync.RWMutex
	clients map[string]*Client
	events  chan statusEvent
}

package main

import (
	"maps"
	"slices"
	"sync"
)

var clientPool = ClientPool{
	clients: make(map[*Client]bool),
}

type ClientPool struct {
	mu      sync.Mutex
	clients map[*Client]bool
}

func (p *ClientPool) addClient(client *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients[client] = true
}

func (p *ClientPool) snapshotClients() []*Client {
	p.mu.Lock()
	defer p.mu.Unlock()

	return slices.Collect(maps.Keys(p.clients))
}

func (p *ClientPool) removeClient(client *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.clients, client)
}

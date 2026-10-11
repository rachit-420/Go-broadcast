package server

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
	_, ok := p.clients[client]
	if ok {
		return
	}
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

func (p *ClientPool) getClientsByRoom(room string) []*Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	clients := make([]*Client, 0)
	for client := range p.clients {
		if client.roomName == room {
			clients = append(clients, client)
		}
	}
	return clients
}

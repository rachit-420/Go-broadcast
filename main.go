package main

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	mu   sync.Mutex
	conn *websocket.Conn
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

func (c *Client) writeMessage(messageType int, message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteMessage(messageType, message)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clientPool = ClientPool{
	clients: make(map[*Client]bool),
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	client := &Client{conn: conn}
	clientPool.addClient(client)
	defer handleClose(client)

	for {
		messageType, message, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		fmt.Println("Received:", string(message))
		broadcast(message, messageType)
	}
}

func broadcast(message []byte, messageType int) {
	clients := clientPool.snapshotClients()

	for _, client := range clients {
		if err := client.writeMessage(messageType, message); err != nil {
			fmt.Println("Write error:", err)
		}
	}
}

func handleClose(client *Client) {
	clientPool.removeClient(client)
	if err := client.conn.Close(); err != nil {
		fmt.Println("Close error:", err)
	}
}

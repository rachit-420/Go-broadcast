package main

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

type Container struct {
	mu   sync.Mutex
	pool map[*Client]bool
}

type Client struct {
	mu   sync.Mutex
	conn *websocket.Conn
}

func (p *Container) inc(conn *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool[conn] = true
}

func (p *Container) lockAndSlice() []*Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	return slices.Collect(maps.Keys(p.pool))
}

func (p *Container) del(conn *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.pool, conn)
}

func (c *Client) writeMessage(message []byte, messageType int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(messageType, message)
}

var upgarder = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connPool = Container{
	pool: make(map[*Client]bool),
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgarder.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &Client{
		conn: conn,
	}
	defer handleClose(client)
	connPool.inc(client)

	for {
		messageType, message, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received", string(message))
		broadCast(message, messageType)
	}
}

func broadCast(message []byte, messageType int) {
	clients := connPool.lockAndSlice()
	for _, conn := range clients {
		err := conn.writeMessage(message, messageType)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func handleClose(conn *Client) {
	connPool.del(conn)
	conn.conn.Close()
}

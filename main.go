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
	pool map[*websocket.Conn]bool
}

func (p *Container) inc(conn *websocket.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool[conn] = true
}

func (p *Container) lockAndSlice() []*websocket.Conn {
	p.mu.Lock()
	defer p.mu.Unlock()
	clients := slices.Collect(maps.Keys(p.pool))
	return clients
}

var upgarder = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connPool = Container{
	pool: make(map[*websocket.Conn]bool),
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
	defer handleClose(conn)
	connPool.inc(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received", string(message))
		broadCast(message, messageType, err)
	}
}

func broadCast(message []byte, messageType int, err error) {
	clients := connPool.lockAndSlice()
	for _, conn := range clients {
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func handleClose(conn *websocket.Conn) {
	delete(connPool.pool, conn)
	conn.Close()
}

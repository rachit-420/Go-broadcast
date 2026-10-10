package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	mu   sync.Mutex
	conn *websocket.Conn
}

func (c *Client) writeMessage(messageType int, message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteMessage(messageType, message)
}

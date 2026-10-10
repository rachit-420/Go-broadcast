package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	mu        sync.Mutex
	conn      *websocket.Conn
	closeOnce sync.Once
}

func (c *Client) writeMessage(messageType int, message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteMessage(messageType, message)
}

func (c *Client) close() {
	c.closeOnce.Do(func() {
		clientPool.removeClient(c)
		if err := c.conn.Close(); err != nil {
			fmt.Println("Close error:", err)
		}
	})
}

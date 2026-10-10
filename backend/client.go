package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	mu        sync.Mutex
	conn      *websocket.Conn
	closeOnce sync.Once
	username  string
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
		mssg := fmt.Sprintf("%s left the chat", c.username)
		rawLeft := Message{
			Type:     "Left",
			Message:  mssg,
			Username: c.username,
		}
		leftMessage, _ := json.Marshal(rawLeft)
		sendNotification(leftMessage)
	})
}

package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
 
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
		releaseUsername(c.username)
		clientPool.removeClient(c)
		if err := c.conn.Close(); err != nil {
			fmt.Println("Close error:", err)
		}
		mssg := fmt.Sprintf("%s left the chat", c.username)
		rawLeft := Message{
			Type:      "Left",
			Message:   mssg,
			Username:  c.username,
			Timestamp: time.Now().UTC(),
		}
		leftMessage, _ := json.Marshal(rawLeft)
		sendNotification(leftMessage)
	})
}

func sendClientHistory(client *Client) {
	pastMessage := history.GetMessage()
	for _, rawMessage := range pastMessage {
		mess, _ := json.Marshal(rawMessage)
		if err := client.writeMessage(1, mess); err != nil {
			fmt.Println("Error in sending history:", err)
			client.close()
		}
	}
}

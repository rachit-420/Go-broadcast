package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("username"))
	if name == "" {
		sendError(w, "Username is required", http.StatusBadRequest)
		return
	}
	if !claimUsername(name) {
		sendError(w, "Username is taken. Please try other name", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		releaseUsername(name)
		return
	}
	client := &Client{conn: conn, username: name}
	clientPool.addClient(client)
	defer client.close()
	sendClientHistory(client)
	mssg := fmt.Sprintf("%s Joined the chat", name)
	rawGreeting := Message{
		Type:      "Joined",
		Message:   mssg,
		Username:  client.username,
		Timestamp: time.Now().UTC(),
	}
	greeting, err := json.Marshal(rawGreeting)
	if err != nil {
		fmt.Println("Error will encoding:", err)
		return
	}
	sendNotification(greeting)

	for {
		messageType, rawMessage, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		res := Message{
			Message:   string(rawMessage),
			Username:  client.username,
			Type:      "Message",
			Timestamp: time.Now().UTC(),
		}
		history.AddMessage(res)
		message, err := json.Marshal(res)
		fmt.Println("Message is", string(message))
		broadcast(message, messageType)
	}
}

func broadcast(message []byte, messageType int) {
	clients := clientPool.snapshotClients()

	for _, client := range clients {
		if err := client.writeMessage(messageType, message); err != nil {
			fmt.Println("Write error:", err)
			client.close()
		}
	}
}

func sendNotification(message []byte) {
	clients := clientPool.snapshotClients()

	for _, client := range clients {
		if err := client.writeMessage(1, message); err != nil {
			fmt.Println("Write error:", err)
			client.close()
		}
	}
}

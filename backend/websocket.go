package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("username")
	if name == "" {
		sendError(w, "Username is required", http.StatusBadRequest)
		return
	}
	_, exists := usernameDB[name]
	if exists {
		sendError(w, "Username is taken. Please try other name", http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	client := &Client{conn: conn, username: name}
	clientPool.addClient(client)
	mssg := fmt.Sprintf("%s Joined the chat", name)
	rawGreeting := Message{
		Type:     "Joined",
		Message:  mssg,
		Username: client.username,
	}
	greeting, err := json.Marshal(rawGreeting)
	sendNotification(greeting)
	defer client.close()

	for {
		messageType, rawMessage, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		res := Message{
			Message:  string(rawMessage),
			Username: client.username,
			Type:     "Message",
		}
		message, err := json.Marshal(res)
		fmt.Println("Message is", string(message))
		broadcast(message, messageType, client.username)
	}
}

func broadcast(message []byte, messageType int, username string) {
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

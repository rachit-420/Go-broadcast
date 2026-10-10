package main

import (
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	client := &Client{conn: conn}
	clientPool.addClient(client)
	defer client.close()

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
			client.close()
		}
	}
}

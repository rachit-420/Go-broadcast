package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	msg "chat-server/internal/message"
	"chat-server/internal/user"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	client, ok := acceptClient(w, r)
	if !ok {
		return
	}
	defer client.close()

	sendClientHistory(client)
	if !announceJoin(client) {
		return
	}
	readMessages(client)
}

func acceptClient(w http.ResponseWriter, r *http.Request) (*Client, bool) {
	name := strings.TrimSpace(r.URL.Query().Get("username"))
	roomName := strings.TrimSpace(r.URL.Query().Get("roomName"))
	if name == "" {
		sendError(w, "Username is required", http.StatusBadRequest)
		return nil, false
	}
	if !user.ClaimUsername(name) {
		sendError(w, "Username is taken. Please try other name", http.StatusBadRequest)
		return nil, false
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		user.ReleaseUsername(name)
		return nil, false
	}
	client := &Client{conn: conn, username: name, roomName: roomName}
	clientPool.addClient(client)
	return client, true
}

func announceJoin(client *Client) bool {
	mssg := fmt.Sprintf("%s Joined the chat", client.username)
	rawGreeting := msg.Message{
		Type:      "Joined",
		Message:   mssg,
		Username:  client.username,
		Timestamp: time.Now().UTC(),
	}
	greeting, err := json.Marshal(rawGreeting)
	if err != nil {
		fmt.Println("Error will encoding:", err)
		return false
	}
	sendNotification(greeting, client.roomName)
	return true
}

func readMessages(client *Client) {
	for {
		messageType, rawMessage, err := client.conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		res := msg.Message{
			Message:   string(rawMessage),
			Username:  client.username,
			Type:      "Message",
			Timestamp: time.Now().UTC(),
		}
		msg.History.AddMessage(res, client.roomName)
		message, err := json.Marshal(res)
		fmt.Println("Message is", string(message))
		broadcast(message, messageType, client.roomName)
	}
}

func broadcast(message []byte, messageType int, room string) {
	clients := clientPool.getClientsByRoom(room)

	for _, client := range clients {
		if err := client.writeMessage(messageType, message); err != nil {
			fmt.Println("Write error:", err)
			client.close()
		}
	}
}

func sendNotification(message []byte, room string) {
	clients := clientPool.snapshotClients()

	for _, client := range clients {
		if client.roomName != room {
			continue
		}
		if err := client.writeMessage(1, message); err != nil {
			fmt.Println("Write error:", err)
			client.close()
		}
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgarder = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var pool = make(map[*websocket.Conn]bool)

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
	pool[conn] = true

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
	for conn := range pool {
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func handleClose(conn *websocket.Conn) {
	delete(pool, conn)
	conn.Close()
}

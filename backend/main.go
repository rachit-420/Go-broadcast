package main

import (
	"fmt"
	"net/http"

	"chat-server/internal/server"
)

func main() {
	http.HandleFunc("/ws", server.HandleWebSocket)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

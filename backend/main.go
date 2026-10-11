package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	usernameMu sync.Mutex
	usernameDB = make(map[string]struct{})
)

func claimUsername(name string) bool {
	usernameMu.Lock()
	defer usernameMu.Unlock()

	if _, exists := usernameDB[name]; exists {
		return false
	}
	usernameDB[name] = struct{}{}
	return true
}

func releaseUsername(name string) {
	usernameMu.Lock()
	defer usernameMu.Unlock()

	delete(usernameDB, name)
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := ErrorResponse{
		Error: message,
	}

	json.NewEncoder(w).Encode(res)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var usernameDB = make(map[string]struct{})

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

package main

type Message struct {
	Type     string `json:"type"`
	Username string `json:"userName"`
	Message  string `json:"message"`
}

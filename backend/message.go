package main

import (
	"sync"
	"time"
)

var history = MessageHistory{
	message: make([]Message, 0, 50),
}

type Message struct {
	Type      string    `json:"type"`
	Username  string    `json:"userName"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type MessageHistory struct {
	mu      sync.Mutex
	message []Message
}

func (m *MessageHistory) AddMessage(message Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.message = append(m.message, message)

	if len(m.message) > 50 {
		m.message = m.message[len(m.message)-50:]
	}
}

func (m *MessageHistory) GetMessage() []Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	messages := make([]Message, len(m.message))
	copy(messages, m.message)
	return messages
}

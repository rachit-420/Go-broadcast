package message

import (
	"sync"
)

var History = MessageHistory{
	message: make(map[string][]Message),
}

type MessageHistory struct {
	mu      sync.Mutex
	message map[string][]Message
}

func (m *MessageHistory) AddMessage(message Message, room string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	messages := append(m.message[room], message)
	if len(messages) > 50 {
		messages = messages[len(messages)-50:]
	}
	m.message[room] = messages
}

func (m *MessageHistory) GetMessage(room string) []Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	messages := make([]Message, len(m.message[room]))
	copy(messages, m.message[room])
	return messages
}

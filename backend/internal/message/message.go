package message

import (
	"time"
)

type Message struct {
	Type      string    `json:"type"`
	Username  string    `json:"userName"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

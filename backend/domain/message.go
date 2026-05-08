package domain

import "time"

// Message Model represents the core business entity
type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageRepository defines the interface for database operations
type MessageRepository interface {
	Store(message *Message) error
	FetchRecent(limit int) ([]Message, error)
}

// PubSub defines the interface for publishing and subscribing to messages
type PubSub interface {
	Publish(channel string, message []byte) error
	Subscribe(channel string, handler func(payload []byte)) error
}

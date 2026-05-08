package usecase

import (
	"encoding/json"
	"backend/domain"
	"time"
)

// MessageUsecase defines the interface for the message usecase
type MessageUsecase interface {
	GetRecentMessages(limit int) ([]domain.Message, error)
	SaveAndPublishMessage(sender, content string) error
	ListenToMessages(handler func(msg domain.Message)) error
}

type messageUsecase struct {
	repo   domain.MessageRepository
	pubsub domain.PubSub
}

// NewMessageUsecase creates a new MessageUsecase
func NewMessageUsecase(r domain.MessageRepository, p domain.PubSub) MessageUsecase {
	return &messageUsecase{
		repo:   r,
		pubsub: p,
	}
}

func (u *messageUsecase) GetRecentMessages(limit int) ([]domain.Message, error) {
	return u.repo.FetchRecent(limit)
}

func (u *messageUsecase) SaveAndPublishMessage(sender, content string) error {
	msg := &domain.Message{
		Sender:    sender,
		Content:   content,
		CreatedAt: time.Now(),
	}

	// 1. Save to DB
	if err := u.repo.Store(msg); err != nil {
		return err
	}

	// 2. Publish to Redis (or other pubsub)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return u.pubsub.Publish("chat_channel", msgBytes)
}

func (u *messageUsecase) ListenToMessages(handler func(msg domain.Message)) error {
	return u.pubsub.Subscribe("chat_channel", func(payload []byte) {
		var msg domain.Message
		if err := json.Unmarshal(payload, &msg); err == nil {
			handler(msg)
		}
	})
}

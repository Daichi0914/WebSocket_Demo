package infrastructure

import (
	"context"
	"backend/domain"

	"github.com/go-redis/redis/v8"
)

type redisPubSub struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisPubSub creates a new redis pubsub client
func NewRedisPubSub(client *redis.Client) domain.PubSub {
	return &redisPubSub{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *redisPubSub) Publish(channel string, message []byte) error {
	return r.client.Publish(r.ctx, channel, message).Err()
}

func (r *redisPubSub) Subscribe(channel string, handler func(payload []byte)) error {
	pubsub := r.client.Subscribe(r.ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		handler([]byte(msg.Payload))
	}
	return nil
}

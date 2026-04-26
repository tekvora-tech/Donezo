package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisBroker implementasi broker pakai Redis pub/sub
type RedisBroker struct {
	client   *redis.Client
	source   string          // Nama service ini (untuk metadata event)
	handlers map[string]HandlerFunc // Map: eventType → handler
}

// HandlerFunc adalah function yang dipanggil saat terima event
type HandlerFunc func(ctx context.Context, event *Event) error

// NewRedisBroker buat instance broker
func NewRedisBroker(redisAddr, source string) *RedisBroker {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &RedisBroker{
		client:   client,
		source:   source,
		handlers: make(map[string]HandlerFunc),
	}
}

// Publish kirim event ke Redis channel
func (b *RedisBroker) Publish(ctx context.Context, channel string, payload any) error {
	event, err := NewEvent(channel, b.source, payload)
	if err != nil {
		return err
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return b.client.Publish(ctx, channel, data).Err()
}

// Subscribe dengar event dari Redis channel
// Blocking — jalankan di goroutine terpisah
func (b *RedisBroker) Subscribe(ctx context.Context, channels ...string) error {
	pubsub := b.client.Subscribe(ctx, channels...)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var event Event
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			fmt.Printf("failed to unmarshal event: %v\n", err)
			continue
		}

		handler, ok := b.handlers[event.Type]
		if !ok {
			fmt.Printf("no handler for event type: %s\n", event.Type)
			continue
		}

		if err := handler(ctx, &event); err != nil {
			fmt.Printf("handler error for %s: %v\n", event.Type, err)
		}
	}

	return nil
}

// RegisterHandler daftarkan handler untuk event type tertentu
func (b *RedisBroker) RegisterHandler(eventType string, fn HandlerFunc) {
	b.handlers[eventType] = fn
}

// Close tutup koneksi Redis
func (b *RedisBroker) Close() error {
	return b.client.Close()
}
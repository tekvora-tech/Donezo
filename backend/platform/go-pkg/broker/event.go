package broker

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event adalah struktur dasar semua event
type Event struct {
	ID        string          `json:"id"`        // UUID unik per event
	Type      string          `json:"type"`      // Nama event: "user.created", "order.paid"
	Source    string          `json:"source"`    // Service yang kirim: "user-service"
	Timestamp int64           `json:"timestamp"` // Unix timestamp
	Payload   json.RawMessage `json:"payload"`   // Data event (bisa apa saja)
}

// NewEvent buat event baru
func NewEvent(eventType, source string, payload any) (*Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:        generateID(),           // UUID random
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now().Unix(),
		Payload:   data,
	}, nil
}

// ParsePayload parse payload ke struct tertentu
func (e *Event) ParsePayload(v any) error {
	return json.Unmarshal(e.Payload, v)
}

func generateID() string {
	// Simplified, pakai UUID library di production
	return uuid.NewString()
}
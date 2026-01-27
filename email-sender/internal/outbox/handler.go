package outbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Handler struct {
	storage   *Storage
	eventType string
}

func NewHandler(storage *Storage, eventType string) *Handler {
	return &Handler{
		storage:   storage,
		eventType: eventType,
	}
}

func (h *Handler) Handle(ctx context.Context, msg kafka.Message) error {
	if !json.Valid(msg.Value) {
		return fmt.Errorf("invalid JSON payload")
	}

	email := OutboxEmail{
		Payload:  msg.Value,
		DedupKey: string(msg.Key),
		Status:   StatusNew,
		Type:     h.eventType,
	}

	if err := h.storage.Save(ctx, email); err != nil {
		return fmt.Errorf("outbox save failed: %w", err)
	}

	return nil
}

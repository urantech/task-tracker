package user

import (
	mykafka "backend-api/pkg/kafka"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		w: mykafka.NewWriter(brokers, "users.registration"),
	}
}

func (p *Producer) ProduceRegistration(ctx context.Context, u UserResponse) error {
	payload, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}

	err = p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(u.Email),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("kafka publish error: %w", err)
	}

	return nil
}

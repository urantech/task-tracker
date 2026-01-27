package user

import (
	"backend-api/internal/kafka_producer"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w kafka_producer.MessageWriter
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		w: kafka_producer.NewWriter(brokers, "users.registration"),
	}
}

func (p *Producer) ProduceRegistration(ctx context.Context, u UserRegisteredEvent) error {
	payload, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}

	log.Printf("publishing registration event: key=%s", u.Email)

	err = p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(u.Email),
		Value: payload,
	})
	if err != nil {
		log.Printf("kafka publish failed: key=%s err=%v", u.Email, err)

		return fmt.Errorf("kafka publish error: %w", err)
	}

	log.Printf("registration event published: key=%s", u.Email)

	return nil
}

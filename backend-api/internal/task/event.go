package task

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
		w: kafka_producer.NewWriter(brokers, "tasks.daily-report"),
	}
}

func (p *Producer) ProduceDailyReport(ctx context.Context, report DailyReportEvent) error {
	payload, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}

	log.Printf("publishing daily report event: key=%d", report.UserID)

	err = p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte{byte(report.UserID)},
		Value: payload,
	})
	if err != nil {
		log.Printf("kafka publish failed: key=%d err=%v", report.UserID, err)

		return fmt.Errorf("kafka publish error: %w", err)
	}

	log.Printf("daily report event published: key=%d", report.UserID)

	return nil
}

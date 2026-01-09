package task

import (
	mykafka "backend-api/pkg/kafka"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w mykafka.MessageWriter
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		w: mykafka.NewWriter(brokers, "tasks.daily-report"),
	}
}

func (p *Producer) ProduceDailyReport(ctx context.Context, report DailyReportMsg) error {
	payload, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("marshal report: %w", err)
	}

	err = p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte{byte(report.UserId)},
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("kafka publish error: %w", err)
	}

	return nil
}

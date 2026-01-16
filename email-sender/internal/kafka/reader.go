package kafka

import (
	"context"
	"email-sender/internal/outbox"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	reader  *kafka.Reader
	handler *outbox.Handler
}

func NewReader(brokers []string, topic, groupID string, msgHandler *outbox.Handler) *Reader {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &Reader{
		reader:  kafkaReader,
		handler: msgHandler,
	}
}

func (r *Reader) Run(ctx context.Context) error {
	defer func() {
		if closeErr := r.reader.Close(); closeErr != nil {
			log.Printf("Error close reader: %v", closeErr)
		}
	}()

	log.Printf(
		"kafka consumer started: topic=%s group=%s brokers=%v",
		r.reader.Config().Topic,
		r.reader.Config().GroupID,
		r.reader.Config().Brokers,
	)

	for {
		msg, err := r.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("consumer context cancelled, stopping")
				return nil
			}

			log.Printf("fetch message error: %v", err)

			return err
		}

		if err := r.handler.Handle(ctx, msg); err != nil {
			log.Printf("handle error: %v", err)
			continue
		}

		log.Printf(
			"message received: topic=%s partition=%d offset=%d key=%s",
			msg.Topic,
			msg.Partition,
			msg.Offset,
			string(msg.Key),
		)

		if err := r.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("error commit message: %w", err)
		}
	}
}

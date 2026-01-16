package outbox

import (
	"context"
	"email-sender/internal/config"
	"email-sender/internal/email"
	"log"
	"time"
)

type Processor struct {
	storage      *Storage
	sender       *email.Sender
	workersCount int
	batchSize    int
	interval     time.Duration
}

func NewProcessor(storage *Storage, sender *email.Sender, cfg *config.Config) *Processor {
	return &Processor{
		storage:      storage,
		sender:       sender,
		workersCount: cfg.OutboxCfg.WorkersCount,
		batchSize:    cfg.OutboxCfg.BatchSize,
		interval:     time.Duration(cfg.OutboxCfg.IntervalSec) * time.Second,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	jobs := make(chan OutboxEmail)

	for i := 0; i < p.workersCount; i++ {
		go p.work(ctx, i, jobs)
	}

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(jobs)
			log.Println("outbox processor stopped")
			return nil
		case <-ticker.C:
			p.enqueue(ctx, jobs)
		}
	}
}

func (p *Processor) work(ctx context.Context, id int, jobs <-chan OutboxEmail) {
	log.Printf("worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d stopped", id)
			return
		case event, ok := <-jobs:
			if !ok {
				return
			}

			p.processOne(ctx, event)
		}
	}
}

func (p *Processor) processOne(ctx context.Context, event OutboxEmail) {
	var err error

	switch event.Type {
	case "WELCOME":
		err = p.sender.SendWelcome(ctx, event.Payload)
	case "REPORT":
		err = p.sender.SendReport(ctx, event.Payload)
	default:
		log.Printf("Unknown event type: %v. Skip processing", event.Type)
		return
	}

	if err != nil {
		attempt := event.Attempts + 1
		next := nextAttempt(attempt)

		err = p.storage.MarkFailed(ctx, event.ID, attempt, next)
		if err != nil {
			log.Printf("mark failed error: %v", err)
		}

		return
	}

	err = p.storage.MarkSent(ctx, event.ID)
	if err != nil {
		log.Printf("mark sent error: %v", err)
	}
}

func (p *Processor) enqueue(ctx context.Context, jobs chan<- OutboxEmail) {
	events, err := p.storage.FetchBatch(ctx, p.batchSize)
	if err != nil {
		log.Printf("outbox fetch error: %v", err)
		return
	}

	for _, event := range events {
		select {
		case jobs <- event:
		case <-ctx.Done():
			return
		}
	}
}

func nextAttempt(attempt int) time.Time {
	backoff := time.Duration(attempt*attempt) * time.Minute
	return time.Now().Add(backoff)
}

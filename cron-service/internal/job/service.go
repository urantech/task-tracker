package job

import (
	"context"
	"fmt"
	"log"

	"github.com/go-co-op/gocron"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) SetupScheduler(ctx context.Context, scheduler *gocron.Scheduler) error {
	configs, err := s.storage.GetAllJobConfigs(ctx)
	if err != nil {
		return fmt.Errorf("setup scheduler: %w", err)
	}

	for _, cfg := range configs {
		_, err := scheduler.
			Every(cfg.Duration).
			Day().
			At(cfg.ExecutionTime).
			Do(func() {
				log.Printf("Executing job: %s", cfg.JobName)
				// TODO: create gRPC client and call the service
			})
		if err != nil {
			return fmt.Errorf("add job %s: %w", cfg.JobName, err)
		}
	}

	return nil
}

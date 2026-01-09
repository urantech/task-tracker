package task

import (
	"backend-api/internal/common"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	storage   *Storage
	validator *validator.Validate
	producer  *Producer
}

func NewService(storage *Storage, producer *Producer) *Service {
	v := common.NewValidator()

	err := v.RegisterValidation("task_status", func(fl validator.FieldLevel) bool {
		status := fl.Field().String()
		switch TaskStatus(status) {
		case StatusTodo, StatusInProgress, StatusDone:
			return true
		default:
			return false
		}
	})
	if err != nil {
		log.Printf("Failed to register validation: %v", err)
	}

	return &Service{
		storage:   storage,
		validator: v,
		producer:  producer,
	}
}

func (s *Service) CreateTask(ctx context.Context, req CreateRequest, userId int64) (Task, error) {
	if err := s.validator.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return Task{}, common.NewValidationError(ve)
		}

		return Task{}, common.ErrInvalidRequest
	}

	task, err := s.storage.Create(ctx, req.Title, req.Description, userId)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

func (s *Service) GetUserTasks(ctx context.Context, userId int64) ([]Task, error) {
	tasks, err := s.storage.GetAllUserTasks(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get all user tasks: %w", err)
	}

	return tasks, nil
}

func (s *Service) UpdateTask(ctx context.Context, taskId int64, userId int64, req UpdateRequest) (Task, error) {
	if req.Title == nil && req.Description == nil && req.Status == nil {
		return Task{}, common.NewValidationErrorFromMap(map[string]string{
			"request": "at least one field must be provided",
		})
	}

	if err := s.validator.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return Task{}, common.NewValidationError(ve)
		}

		return Task{}, common.ErrInvalidRequest
	}

	task, err := s.storage.UpdateTask(ctx, taskId, userId, req)
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}

func (s *Service) CollectAndSendTaskAnalytics(ctx context.Context) error {
	reports, err := s.storage.GetDailyReports(ctx)
	if err != nil {
		return common.ErrCollectingDailyReports
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(reports))

	for _, report := range reports {
		if report.PendingCount <= 0 && report.CompletedCount <= 0 {
			continue
		}

		wg.Add(1)
		go func(r DailyReport) {
			defer wg.Done()

			reportMsg := DailyReportMsg{
				UserId:  r.UserId,
				Message: buildMessage(r.PendingCount, r.CompletedCount),
			}

			if err := s.produceEvent(ctx, reportMsg); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(report)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return fmt.Errorf("failed to send %d reports", len(errCh))
	}

	return nil
}

func buildMessage(pending, completed int) string {
	switch {
	case pending > 0 && completed > 0:
		return fmt.Sprintf("У вас осталось %d несделанных задач. За сегодня вы выполнили %d задач.",
			pending, completed)
	case pending > 0:
		return fmt.Sprintf("У вас осталось %d несделанных задач.", pending)
	case completed > 0:
		return fmt.Sprintf("За сегодня вы выполнили %d задач.", completed)
	default:
		return ""
	}
}

func (s *Service) produceEvent(ctx context.Context, report DailyReportMsg) error {
	if err := s.producer.ProduceDailyReport(ctx, report); err != nil {
		log.Printf("ERROR: failed to publish daily report for user %d: %v",
			report.UserId, err)
		return err
	}

	return nil
}

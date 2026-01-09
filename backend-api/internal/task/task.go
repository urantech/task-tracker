package task

import "time"

type TaskStatus string

const (
	StatusTodo       TaskStatus = "TODO"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
)

type Task struct {
	Id          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	UserId      int64      `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateRequest struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"max=100"`
}

type UpdateRequest struct {
	Title       *string     `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string     `json:"description,omitempty" validate:"omitempty,max=100"`
	Status      *TaskStatus `json:"status,omitempty" validate:"omitempty,task_status"`
}

type DailyReport struct {
	UserId         int64 `json:"user_id"`
	PendingCount   int   `json:"pending_count"`
	CompletedCount int   `json:"completed_count"`
}

type DailyReportMsg struct {
	UserId  int64  `json:"user_id"`
	Message string `json:"message"`
}

package task

import "time"

type BaseTask struct {
	Name       string    `json:"name" binding:"required"`
	Payload    string    `json:"payload"`
	RunAt      time.Time `json:"run_at" binding:"omitempty"`
	CreatedAt  time.Time `json:"created_at" binding:"omitempty"`
	MaxRetries int       `json:"max_retries" binding:"omitempty"`
	Retries    int       `json:"retries" binding:"omitempty"`
	NextRunAt  time.Time `json:"next_run_at" binding:"omitempty"`
}

type CreateTaskRequest struct {
	BaseTask
	State string `json:"state" binding:"omitempty"`
}

type HandleTask struct {
	ID    string `json:"id" binding:"required"`
	State string `json:"state" binding:"required"`
	BaseTask
}

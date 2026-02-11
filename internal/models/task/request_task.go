package task

import "time"

type RequestTask struct {
	Name       string    `json:"name" binding:"required"`
	Payload    string    `json:"payload"`
	RunAt      time.Time `json:"run_at" binding:"omitempty"`
	State      string    `json:"state" binding:"omitempty"`
	CreatedAt  time.Time `json:"created_at" binding:"omitempty"`
	MaxRetries int       `json:"max_retries" binding:"omitempty"`
	Retries    int       `json:"retries" binding:"omitempty"`
	NextRunAt  time.Time `json:"next_run_at" binding:"omitempty"`
}

type RequestTasks struct {
	ID         string    `json:"id" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Payload    string    `json:"payload"`
	RunAt      time.Time `json:"run_at" binding:"omitempty"`
	State      string    `json:"state" binding:"required"`
	CreatedAt  time.Time `json:"created_at" binding:"omitempty"`
	MaxRetries int       `json:"max_retries" binding:"omitempty"`
	Retries    int       `json:"retries" binding:"omitempty"`
	NextRunAt  time.Time `json:"next_run_at" binding:"omitempty"`
}

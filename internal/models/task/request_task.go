package task

import "time"

type RequestTask struct {
	Name    string    `json:"name" binding:"required"`
	Payload string    `json:"payload"`
	RunAt   time.Time `json:"run_at" binding:"omitempty"`
}

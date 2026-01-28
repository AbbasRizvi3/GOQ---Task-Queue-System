package task

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskState string

const (
	TaskPending   TaskState = "pending"
	TaskLeased    TaskState = "leased"
	TaskRetry     TaskState = "retry"
	TaskCompleted TaskState = "completed"
	TaskFailed    TaskState = "failed"
	TaskDead      TaskState = "dead"
)

type Task struct {
	ID          string
	Name        string
	State       TaskState
	Payload     []byte
	Attempts    int
	MaxAttempts int
	RunAt       time.Time
	LeaseUntil  time.Time
	LastError   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(name string, payload []byte, maxAttempts int, runAt time.Time) *Task {
	return &Task{
		ID:          fmt.Sprintf("%v", uuid.New().String()),
		Name:        name,
		State:       TaskPending,
		Payload:     payload,
		Attempts:    0,
		MaxAttempts: maxAttempts,
		RunAt:       runAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (t *Task) IsReadyToRun() bool {
	if time.Now().After(t.RunAt) && (t.State == TaskPending || t.State == TaskRetry) {
		return true
	}
	return false
}

func (t *Task) CanRetry() bool {
	return t.Attempts < t.MaxAttempts
}

func (t *Task) MarkLeased(leaseDuration time.Duration) {
	t.State = TaskLeased
	t.LeaseUntil = time.Now().Add(leaseDuration)
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkCompleted() {
	t.State = TaskCompleted
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkFailed(errMsg string) {
	t.Attempts++
	t.LastError = errMsg
	t.UpdatedAt = time.Now()
	if t.CanRetry() {
		t.State = TaskRetry
	} else {
		t.State = TaskDead
	}
}

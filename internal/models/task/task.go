package task

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	maxRetriesDefault = 3
)

type Task struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	Payload    string    `json:"payload"`
	State      string    `json:"state"`
	RunAt      time.Time `json:"run_at"`
	CreatedAt  time.Time `json:"created_at"`
	NextRunAt  time.Time `json:"next_run_at"`
	LeaseUntil time.Time `json:"lease_until"`
	MaxRetries int
	Retries    int
	Error      string
	UpdatedAt  time.Time
	Mu         sync.Mutex `json:"-"`
}

func NewTask(name string, payload string, runAt time.Time) *Task {
	return &Task{
		ID:         uuid.New().String()[:8],
		Name:       name,
		Payload:    payload,
		State:      "pending",
		RunAt:      runAt,
		MaxRetries: maxRetriesDefault,
		Retries:    0,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}
}

func (t *Task) IsReadyToRun() bool {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if (t.State == "pending" || t.State == "retry") && (t.RunAt.Before(time.Now()) || t.NextRunAt.Before(time.Now()) || t.NextRunAt.Equal(time.Now()) || t.RunAt.Equal(time.Now())) {
		return true
	}
	return false
}

func (t *Task) MarkLeased(leaseSecs int) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.State = "leased"
	t.LeaseUntil = time.Now().Add(time.Duration(leaseSecs) * time.Second)
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkCompleted() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.State = "completed"
	t.UpdatedAt = time.Now()
	t.NextRunAt = time.Time{}
}

func (t *Task) MarkFailed(reason string) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.Retries++
	if t.Retries > t.MaxRetries {
		t.State = "dead"
		t.Error = fmt.Sprintf("Max retries reached (%d/%d). Last error: %s", t.Retries-1, t.MaxRetries, reason)
		t.UpdatedAt = time.Now()
		t.NextRunAt = time.Time{}
		fmt.Printf("Fatal: Task %s is DEAD. Max retries (%d) exceeded. Error: %s\n", t.ID, t.MaxRetries, reason)
		return
	}
	now := time.Now()
	t.State = "retry"
	t.Error = reason
	t.UpdatedAt = now
	shift := uint(0)
	if t.Retries > 1 {
		shift = uint(t.Retries - 1)
	}
	baseDelay := 60 * time.Second
	backoff := baseDelay * time.Duration(1<<shift)
	halfBackoff := int64(backoff) / 2
	var jitter time.Duration
	if halfBackoff > 0 {
		jitter = time.Duration(rand.Int63n(halfBackoff)) - backoff/4
	}

	t.NextRunAt = now.Add(backoff + jitter)

	fmt.Printf("Retryable: Task %s failed, scheduled retry (%d/%d) at %v\n",
		t.ID, t.Retries, t.MaxRetries, t.NextRunAt)
}

func (t *Task) MarkDead(reason string) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.State = "dead"
	t.Error = reason
	t.UpdatedAt = time.Now()
	t.NextRunAt = time.Time{}
	fmt.Printf("Fatal: Task %s is DEAD. Max retries (%d) exceeded. Error: %s\n", t.ID, t.MaxRetries, reason)
}

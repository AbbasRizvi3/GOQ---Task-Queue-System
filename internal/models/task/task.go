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
	if (t.State == "pending" || t.State == "retry" || t.State == "ready") && (t.RunAt.Before(time.Now()) || t.NextRunAt.Before(time.Now()) || t.NextRunAt.Equal(time.Now()) || t.RunAt.Equal(time.Now())) {
		return true
	}
	return false
}

func (t *Task) GetState() string {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.State
}

func (t *Task) GetNextRunAt() time.Time {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.NextRunAt
}

func (t *Task) GetRetries() int {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.Retries
}

func (t *Task) GetID() string {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.ID
}

func (t *Task) GetRunAt() time.Time {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.RunAt
}
func (t *Task) GetLeaseUntil() time.Time {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.LeaseUntil
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

func (t *Task) MarkCanceled() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.Retries++
	isMaxed := t.Retries >= t.MaxRetries
	if isMaxed {
		t.MarkDead("Retries Exceeded")
	} else {
		t.State = "canceled"
		t.UpdatedAt = time.Now()
		t.NextRunAt = time.Time{}
	}
}

func (t *Task) MarkReady() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.State = "ready"
	t.UpdatedAt = time.Now()
	t.NextRunAt = time.Now()
}

func (t *Task) MarkFailed(reason string) {
	t.Mu.Lock()
	t.Retries++
	isMaxed := t.Retries >= t.MaxRetries
	t.Mu.Unlock()

	if isMaxed {
		t.MarkDead(reason)
	} else {
		now := time.Now()
		t.Mu.Lock()
		t.State = "retry"
		baseDelay := 60 * time.Second
		backoff := baseDelay * time.Duration(1<<uint(t.Retries-1))
		jitter := time.Duration(rand.Int63n(int64(backoff)/2)) - backoff/4

		t.NextRunAt = now.Add(backoff + jitter)
		t.Error = reason
		t.UpdatedAt = now
		t.Mu.Unlock()
		fmt.Printf("Retryable: Task %s failed, will retry (attempt %d/%d) at %v (in ~%v)\n", t.ID, t.Retries, t.MaxRetries, t.NextRunAt, backoff+jitter)
	}
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

func (t *Task) MarkRetry() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.State = "retry"
	t.UpdatedAt = time.Now()
	t.NextRunAt = time.Now()
}

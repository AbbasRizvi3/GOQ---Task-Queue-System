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
	ID         string
	Payload    []byte
	State      string
	RunAt      time.Time
	NextRunAt  time.Time
	LeaseUntil time.Time
	MaxRetries int
	Retries    int
	Error      string
	UpdatedAt  time.Time
	Mu         sync.Mutex
}

func NewTask(name string, payload []byte, runAt time.Time) *Task {
	return &Task{
		ID:         uuid.New().String()[:8],
		Payload:    payload,
		State:      "pending",
		RunAt:      runAt,
		MaxRetries: maxRetriesDefault,
		Retries:    0,
		UpdatedAt:  time.Now(),
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

func (t *Task) GetState() string {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.State
}

func (t *Task) GetLeaseUntil() time.Time {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.LeaseUntil
}

func (t *Task) MarkReady() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.State = "ready"
	t.UpdatedAt = time.Now()
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
}

func (t *Task) MarkFailed(reason string) {
	t.Mu.Lock()
	t.Retries++
	isMaxed := t.Retries >= t.MaxRetries
	t.Mu.Unlock()

	if isMaxed {
		t.MarkDead(reason)
	} else {
		t.Mu.Lock()
		t.State = "retry"

		baseDelay := 2 * time.Second
		backoff := baseDelay * (1 << (t.Retries - 1))
		jitter := time.Duration(rand.Int63n(int64(backoff)/2)) - backoff/4

		t.NextRunAt = time.Now().Add(backoff + jitter)
		t.Error = reason
		t.UpdatedAt = time.Now()
		t.Mu.Unlock()
		fmt.Printf("Retryable: Task %s failed, will retry (attempt %d/%d)\n", t.ID, t.Retries, t.MaxRetries)
	}
}

func (t *Task) MarkDead(reason string) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.State = "dead"
	t.Error = reason
	t.UpdatedAt = time.Now()
	fmt.Printf("Fatal: Task %s is DEAD. Max retries (%d) exceeded. Error: %s\n", t.ID, t.MaxRetries, reason)
}

func (t *Task) MarkRetry() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.State = "retry"
	t.UpdatedAt = time.Now()
	t.NextRunAt = time.Now()
}

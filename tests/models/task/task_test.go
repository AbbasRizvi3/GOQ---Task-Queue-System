package task

import (
	"testing"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
)

func TestNewTask(t *testing.T) {
	name := "test task"
	payload := "test payload"
	runAt := time.Now()

	tsk := task.NewTask(name, payload, runAt)

	if tsk.Name != name {
		t.Errorf("Expected name %s, got %s", name, tsk.Name)
	}
	if tsk.Payload != payload {
		t.Errorf("Expected payload %s, got %s", payload, tsk.Payload)
	}
	if tsk.State != "pending" {
		t.Errorf("Expected initial state 'pending', got %s", tsk.State)
	}
	if tsk.ID == "" {
		t.Error("Expected ID to be set")
	}
	if tsk.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", tsk.MaxRetries)
	}
}

func TestGetState(t *testing.T) {
	tsk := task.NewTask("test", "payload", time.Now())
	state := tsk.GetState()

	if state != "pending" {
		t.Errorf("Expected state 'pending', got %s", state)
	}
}

func TestIsReadyToRun(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)
	futureTime := now.Add(1 * time.Hour)

	tests := []struct {
		name      string
		state     string
		runAt     time.Time
		nextRunAt time.Time
		expected  bool
	}{
		{
			name:      "ready_pending",
			state:     "pending",
			runAt:     pastTime,
			nextRunAt: time.Time{},
			expected:  true,
		},
		{
			name:      "ready_retry",
			state:     "retry",
			runAt:     now,
			nextRunAt: now.Add(-1 * time.Second),
			expected:  true,
		},
		{
			name:      "not_ready_future",
			state:     "pending",
			runAt:     futureTime,
			nextRunAt: futureTime,
			expected:  false,
		},
		{
			name:      "ready_ready_state",
			state:     "ready",
			runAt:     time.Time{},
			nextRunAt: now.Add(-1 * time.Second),
			expected:  true,
		},
		{
			name:      "not_ready_leased",
			state:     "leased",
			runAt:     pastTime,
			nextRunAt: time.Time{},
			expected:  false,
		},
		{
			name:      "not_ready_completed",
			state:     "completed",
			runAt:     pastTime,
			nextRunAt: time.Time{},
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsk := task.NewTask("test", "payload", tt.runAt)
			tsk.State = tt.state
			tsk.NextRunAt = tt.nextRunAt

			result := tsk.IsReadyToRun()
			if result != tt.expected {
				t.Errorf("Expected IsReadyToRun to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMarkLeased(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.MarkLeased(15)

	if tsk.GetState() != "leased" {
		t.Errorf("Expected state 'leased', got %s", tsk.GetState())
	}

	leaseUntil := tsk.GetLeaseUntil()
	if leaseUntil.IsZero() {
		t.Error("LeaseUntil should not be zero after marking leased")
	}

	diff := time.Until(leaseUntil)
	if diff < 14*time.Second || diff > 16*time.Second {
		t.Errorf("Expected lease to be ~15 seconds, got %v", diff)
	}
}

func TestMarkCompleted(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.MarkCompleted()

	if tsk.GetState() != "completed" {
		t.Errorf("Expected state 'completed', got %s", tsk.GetState())
	}
}

func TestMarkCanceled(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.NextRunAt = time.Now().Add(1 * time.Hour)

	tsk.MarkCanceled()

	if tsk.GetState() != "canceled" {
		t.Errorf("Expected state 'canceled', got %s", tsk.GetState())
	}
}

func TestMarkReady(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.State = "pending"

	if tsk.GetState() != "pending" {
		t.Errorf("Expected state 'pending', got %s", tsk.GetState())
	}

	if tsk.NextRunAt.IsZero() {
		t.Error("NextRunAt should be set after MarkReady")
	}
}

func TestMarkFailed_Retryable(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.MaxRetries = 3
	tsk.Retries = 0

	tsk.MarkFailed("test error")

	if tsk.GetState() != "retry" {
		t.Errorf("Expected state 'retry', got %s", tsk.GetState())
	}
	if tsk.Retries != 1 {
		t.Errorf("Expected retries to be 1, got %d", tsk.Retries)
	}
	if tsk.Error != "test error" {
		t.Errorf("Expected error 'test error', got %s", tsk.Error)
	}
}

func TestMarkFailed_Dead(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.MaxRetries = 2
	tsk.Retries = 2

	tsk.MarkFailed("fatal error")

	if tsk.GetState() != "dead" {
		t.Errorf("Expected state 'dead', got %s", tsk.GetState())
	}
	if tsk.Error != "fatal error" {
		t.Errorf("Expected error 'fatal error', got %s", tsk.Error)
	}
}

func TestMarkDead(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.MarkDead("test reason")

	if tsk.GetState() != "dead" {
		t.Errorf("Expected state 'dead', got %s", tsk.GetState())
	}
	if tsk.Error != "test reason" {
		t.Errorf("Expected error 'test reason', got %s", tsk.Error)
	}
}

func TestMarkRetry(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.State = "retry"

	if tsk.GetState() != "retry" {
		t.Errorf("Expected state 'retry', got %s", tsk.GetState())
	}
	if tsk.NextRunAt.IsZero() {
		t.Error("NextRunAt should be set after MarkRetry")
	}
}

func TestConcurrentStateChanges(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			tsk.MarkLeased(5)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			tsk.State = "ready"
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			_ = tsk.GetState()
		}
		done <- true
	}()

	<-done
	<-done
	<-done

	state := tsk.GetState()
	if state == "" {
		t.Error("State should not be empty")
	}
}

func TestRetryBackoff(t *testing.T) {
	tsk := task.NewTask("test task", "payload", time.Now())
	tsk.MaxRetries = 5

	for i := 0; i < 4; i++ {
		tsk.MarkFailed("error")
		state := tsk.GetState()
		if state != "retry" {
			t.Errorf("Iteration %d: Expected state 'retry', got %s", i, state)
		}
		if tsk.Retries != i+1 {
			t.Errorf("Iteration %d: Expected retries %d, got %d", i, i+1, tsk.Retries)
		}
	}

	tsk.MarkFailed("error")
	if tsk.GetState() != "dead" {
		t.Errorf("Expected state 'dead' after max retries, got %s", tsk.GetState())
	}
}

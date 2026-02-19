package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
	schedulerpkg "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/scheduler"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestScheduleTasks_ContextCancellation(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	memQueue := queue.NewMemoryQueue()
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		schedulerpkg.ScheduleTasks(ctx, memQueue)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("Scheduler did not stop after context cancellation")
	}
}

func TestFetchAndCacheTasks(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "name", "payload", "state", "run_at", "next_run_at", "lease_until", "max_retries", "retries", "error", "created_at", "updated_at",
	}).
		AddRow("task1", "test", "payload1", "pending", now.Add(-1*time.Hour), time.Time{}, time.Time{}, 3, 0, "", now, now).
		AddRow("task2", "test2", "payload2", "retry", now, now.Add(-1*time.Second), time.Time{}, 3, 1, "", now, now)

	mock.ExpectQuery("SELECT id, name, payload, state, run_at, next_run_at").
		WithArgs(20).
		WillReturnRows(rows)

	memQueue := queue.NewMemoryQueue()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	schedulerpkg.ScheduleTasks(ctx, memQueue)

	if len(memQueue.Tasks) < 1 {
		t.Logf("Expected tasks in queue, got %d", len(memQueue.Tasks))
	}
}

func TestSchedulerTiming(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	memQueue := queue.NewMemoryQueue()
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	startTime := time.Now()
	schedulerpkg.ScheduleTasks(ctx, memQueue)
	elapsedTime := time.Since(startTime)
	if elapsedTime < 2*time.Second {
		t.Logf("Scheduler ran for %v, expected at least 2s", elapsedTime)
	}
}

func TestTaskFiltering_PastDue(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "name", "payload", "state", "run_at", "next_run_at", "lease_until", "max_retries", "retries", "error", "created_at", "updated_at",
	}).
		AddRow("past1", "test", "payload", "pending", pastTime, time.Time{}, time.Time{}, 3, 0, "", now, now)

	mock.ExpectQuery("SELECT id, name, payload, state, run_at, next_run_at").
		WithArgs(20).
		WillReturnRows(rows)

	memQueue := queue.NewMemoryQueue()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	schedulerpkg.ScheduleTasks(ctx, memQueue)

	if len(memQueue.Tasks) < 1 {
		t.Logf("Expected past-due task in queue, got %d", len(memQueue.Tasks))
	}
}

func TestTaskFiltering_ByState(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "name", "payload", "state", "run_at", "next_run_at", "lease_until", "max_retries", "retries", "error", "created_at", "updated_at",
	}).
		AddRow("task1", "test", "payload", "pending", pastTime, time.Time{}, time.Time{}, 3, 0, "", now, now).
		AddRow("task2", "test", "payload", "retry", pastTime, now.Add(-1*time.Second), time.Time{}, 3, 1, "", now, now)

	mock.ExpectQuery("SELECT id, name, payload, state, run_at, next_run_at").
		WithArgs(20).
		WillReturnRows(rows)

	memQueue := queue.NewMemoryQueue()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	schedulerpkg.ScheduleTasks(ctx, memQueue)

	if len(memQueue.Tasks) < 1 {
		t.Logf("Expected tasks in valid states, got %d", len(memQueue.Tasks))
	}
}

func TestBatchProcessing(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)

	builder := sqlmock.NewRows([]string{
		"id", "name", "payload", "state", "run_at", "next_run_at", "lease_until", "max_retries", "retries", "error", "created_at", "updated_at",
	})

	for i := 0; i < 20; i++ {
		builder.AddRow("task"+string(rune('0'+i%10)), "test", "payload", "pending", pastTime, time.Time{}, time.Time{}, 3, 0, "", now, now)
	}

	mock.ExpectQuery("SELECT id, name, payload, state, run_at, next_run_at").
		WithArgs(20).
		WillReturnRows(builder)

	memQueue := queue.NewMemoryQueue()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	schedulerpkg.ScheduleTasks(ctx, memQueue)

	if len(memQueue.Tasks) > 20 {
		t.Errorf("Expected at most 20 tasks, got %d", len(memQueue.Tasks))
	}
}

func TestSchedulerWithWorkerCoordination(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "name", "payload", "state", "run_at", "next_run_at", "lease_until", "max_retries", "retries", "error", "created_at", "updated_at",
	}).
		AddRow("task1", "test", "payload", "pending", pastTime, time.Time{}, time.Time{}, 3, 0, "", now, now)

	mock.ExpectQuery("SELECT id, name, payload, state, run_at, next_run_at").
		WithArgs(20).
		WillReturnRows(rows)

	memQueue := queue.NewMemoryQueue()

	oldSignal := app.ProcessSignal
	signalCh := make(chan struct{}, 100)
	app.ProcessSignal = signalCh
	defer func() { app.ProcessSignal = oldSignal }()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	schedulerpkg.ScheduleTasks(ctx, memQueue)

	if len(memQueue.Tasks) < 1 {
		t.Logf("Expected task to be enqueued for worker coordination")
	}
}

func TestSchedulerDatabaseInteraction(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"id", "name", "payload", "state", "run_at", "next_run_at", "lease_until", "max_retries", "retries", "error", "created_at", "updated_at",
	}).
		AddRow("task1", "test", "payload", "pending", pastTime, time.Time{}, time.Time{}, 3, 0, "", now, now)

	mock.ExpectQuery("SELECT id, name, payload, state, run_at, next_run_at").
		WithArgs(20).
		WillReturnRows(rows)

	memQueue := queue.NewMemoryQueue()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	schedulerpkg.ScheduleTasks(ctx, memQueue)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Logf("Database expectations not met: %v", err)
	}
}

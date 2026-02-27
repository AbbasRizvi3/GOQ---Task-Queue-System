package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	workerpkg "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/worker"
	"github.com/DATA-DOG/go-sqlmock"
)

func setupTestCtx() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func TestProcessTask_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()
	app.Databasehandle = mockDB
	app.ProcessSignal = make(chan *task.Task, 100)

	ctx, cancel := setupTestCtx()
	defer cancel()
	now := time.Now()
	leaseTime := now.Add(8 * time.Second)
	tsk := &task.Task{
		ID:         "task-1",
		State:      "leased",
		LeaseUntil: leaseTime,
	}

	mock.ExpectQuery("SELECT state, lease_until FROM tasks WHERE id =").
		WithArgs(tsk.ID).
		WillReturnRows(sqlmock.NewRows([]string{"state", "lease_until"}).
			AddRow("leased", leaseTime))

	mock.ExpectExec("^UPDATE tasks SET state").WillReturnResult(sqlmock.NewResult(0, 1))

	workerpkg.ProcessTask(ctx)
	app.ProcessSignal <- tsk

	time.Sleep(6 * time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %v", err)
	}
}

func TestTaskCancellation_WorkerAborts(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	app.Databasehandle = mockDB
	app.ProcessSignal = make(chan *task.Task, 100)

	ctx, cancel := setupTestCtx()
	defer cancel()

	leaseTime := time.Now().Add(8 * time.Second)
	tsk := &task.Task{
		ID:         "canceled-task",
		State:      "leased",
		LeaseUntil: leaseTime,
	}

	mock.ExpectQuery("SELECT state, lease_until FROM tasks WHERE id =").
		WithArgs(tsk.ID).
		WillReturnRows(sqlmock.NewRows([]string{"state", "lease_until"}).
			AddRow("canceled", leaseTime))

	workerpkg.ProcessTask(ctx)
	app.ProcessSignal <- tsk

	time.Sleep(6 * time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Worker did not abort as expected: %v", err)
	}
}

func TestWorkerOwnershipCheck_LeaseChanged(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()
	app.Databasehandle = mockDB
	app.ProcessSignal = make(chan *task.Task, 100)

	ctx, cancel := setupTestCtx()
	defer cancel()

	originalLease := time.Now().Add(8 * time.Second)
	newLease := time.Now().Add(16 * time.Second)

	tsk := &task.Task{
		ID:         "stolen-task",
		State:      "leased",
		LeaseUntil: originalLease,
	}

	mock.ExpectQuery("SELECT state, lease_until FROM tasks WHERE id =").
		WithArgs(tsk.ID).
		WillReturnRows(sqlmock.NewRows([]string{"state", "lease_until"}).
			AddRow("leased", newLease))

	workerpkg.ProcessTask(ctx)
	app.ProcessSignal <- tsk

	time.Sleep(6 * time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Worker should have detected lease mismatch: %v", err)
	}
}

func TestSyncTaskToDB(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()
	app.Databasehandle = mockDB

	tsk := &task.Task{
		ID:      "sync-id",
		State:   "completed",
		Retries: 1,
	}

	mock.ExpectExec("^UPDATE tasks SET state").
		WithArgs("completed", nil, tsk.LeaseUntil, 1, "", sqlmock.AnyArg(), tsk.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	workerpkg.SyncTaskToDB(tsk)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SyncTaskToDB failed: %v", err)
	}
}

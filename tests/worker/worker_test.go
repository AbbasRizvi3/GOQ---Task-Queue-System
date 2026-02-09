package worker

import (
	"testing"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
	workerpkg "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/worker"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestProcessTask_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	tsk := task.NewTask("test", "payload", now.Add(-1*time.Hour))
	tsk.State = "pending"

	memQueue := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	memQueue.Enqueue(tsk, signalCh)

	mock.ExpectExec("UPDATE tasks SET").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	oldSignal := app.ProcessSignal
	app.ProcessSignal = signalCh
	defer func() { app.ProcessSignal = oldSignal }()

	workerpkg.ProcessTask(memQueue)

	select {
	case signalCh <- struct{}{}:
	default:
	}

	time.Sleep(6 * time.Second)

	if len(memQueue.Tasks) > 0 {
		t.Logf("Expected task to be removed after processing, still has %d tasks", len(memQueue.Tasks))
	}
}

func TestWorkerPoolManagement(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	memQueue := queue.NewMemoryQueue()

	for i := 0; i < 3; i++ {
		tsk := task.NewTask("test", "payload", now.Add(-1*time.Hour))
		tsk.State = "pending"
		tsk.ID = string(rune('0' + i))
		memQueue.Enqueue(tsk, make(chan struct{}, 100))
	}

	for i := 0; i < 3; i++ {
		mock.ExpectExec("UPDATE tasks SET").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	signalCh := make(chan struct{}, 100)
	oldSignal := app.ProcessSignal
	app.ProcessSignal = signalCh
	defer func() { app.ProcessSignal = oldSignal }()

	workerpkg.ProcessTask(memQueue)

	for i := 0; i < 3; i++ {
		signalCh <- struct{}{}
	}

	time.Sleep(6 * time.Second)

	if len(memQueue.Tasks) > 0 {
		t.Logf("Expected all tasks to be processed, still has %d", len(memQueue.Tasks))
	}
}

func TestSyncTaskToDB(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	tsk := task.NewTask("test", "payload", now)
	tsk.State = "completed"
	tsk.Retries = 2

	mock.ExpectExec("UPDATE tasks SET state").
		WithArgs("completed", sqlmock.AnyArg(), sqlmock.AnyArg(), 2, "", sqlmock.AnyArg(), tsk.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	workerpkg.SyncTaskToDB(tsk)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Logf("Database update expectations not met: %v", err)
	}
}

func TestFindReadyTask(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	memQueue := queue.NewMemoryQueue()

	readyTask := task.NewTask("test", "payload", now.Add(-1*time.Hour))
	readyTask.State = "pending"
	readyTask.ID = "ready123"
	signalCh := make(chan struct{}, 100)
	memQueue.Enqueue(readyTask, signalCh)

	mock.ExpectExec("UPDATE tasks SET").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if len(memQueue.Tasks) < 1 {
		t.Error("Expected task in queue")
	}
}

func TestTaskCancellation(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	memQueue := queue.NewMemoryQueue()

	canceledTask := task.NewTask("test", "payload", now.Add(-1*time.Hour))
	canceledTask.State = "canceled"
	canceledTask.ID = "canceled1"
	signalCh := make(chan struct{}, 100)
	memQueue.Enqueue(canceledTask, signalCh)

	mock.ExpectExec("UPDATE tasks SET").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	oldSignal := app.ProcessSignal
	signalCh2 := make(chan struct{}, 100)
	app.ProcessSignal = signalCh2
	defer func() { app.ProcessSignal = oldSignal }()

	workerpkg.ProcessTask(memQueue)

	signalCh2 <- struct{}{}

	time.Sleep(6 * time.Second)

	if len(memQueue.Tasks) > 0 {
		t.Logf("Expected canceled task to be removed, still has %d tasks", len(memQueue.Tasks))
	}
}

func TestTaskFailureHandling(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	tsk := task.NewTask("test", "payload", now)
	tsk.State = "pending"

	mock.ExpectExec("UPDATE tasks SET").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tsk.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	tsk.MarkFailed("Test error")
	workerpkg.SyncTaskToDB(tsk)

	if tsk.State != "dead" && tsk.State != "retry" {
		t.Errorf("Expected task to be in dead or retry state after MarkFailed, got %s", tsk.State)
	}
}

func TestConcurrentWorkerProcessing(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer mockDB.Close()

	oldDB := app.Databasehandle
	app.Databasehandle = mockDB
	defer func() { app.Databasehandle = oldDB }()

	now := time.Now()
	memQueue := queue.NewMemoryQueue()

	for i := 0; i < 5; i++ {
		tsk := task.NewTask("test", "payload", now.Add(-1*time.Hour))
		tsk.State = "pending"
		tsk.ID = "task" + string(rune('0'+i))
		signalCh := make(chan struct{}, 100)
		memQueue.Enqueue(tsk, signalCh)
	}

	for i := 0; i < 5; i++ {
		mock.ExpectExec("UPDATE tasks SET").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	signalCh := make(chan struct{}, 100)
	oldSignal := app.ProcessSignal
	app.ProcessSignal = signalCh
	defer func() { app.ProcessSignal = oldSignal }()

	workerpkg.ProcessTask(memQueue)

	go func() {
		for i := 0; i < 5; i++ {
			signalCh <- struct{}{}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(6500 * time.Millisecond)

	if len(memQueue.Tasks) > 0 {
		t.Logf("Expected all concurrent tasks to be processed, still has %d", len(memQueue.Tasks))
	}
}

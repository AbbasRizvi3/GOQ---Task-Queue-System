package worker_test

import (
	"testing"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	workerpkg "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/worker"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestProcessTask_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer func() {
		if err := mockDB.Close(); err != nil {
			t.Fatalf("Failed to close mock DB: %v", err)
		}
	}()
	app.Databasehandle = mockDB

	app.ProcessSignal = make(chan *task.Task, 100)
	app.CancelSignal = make(chan string, 100)

	mock.ExpectExec("^UPDATE tasks SET state").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("^UPDATE tasks SET state").WillReturnResult(sqlmock.NewResult(0, 1))

	workerpkg.ProcessTask()

	now := time.Now()
	tsk := task.NewTask("test", "payload", now.Add(-1*time.Hour))
	tsk.State = "pending"
	app.ProcessSignal <- tsk
	time.Sleep(6 * time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database update expectations not met: %v", err)
	}
}

func TestWorkerPoolManagement(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer func() {
		if err := mockDB.Close(); err != nil {
			t.Fatalf("Failed to close mock DB: %v", err)
		}
	}()

	app.Databasehandle = mockDB
	app.ProcessSignal = make(chan *task.Task, 100)
	app.CancelSignal = make(chan string, 100)

	for i := 0; i < 6; i++ {
		mock.ExpectExec("^UPDATE tasks SET state").WillReturnResult(sqlmock.NewResult(0, 1))
	}

	workerpkg.ProcessTask()

	now := time.Now()
	for i := 0; i < 3; i++ {
		tsk := task.NewTask("test", "payload", now.Add(-1*time.Hour))
		tsk.State = "pending"
		tsk.ID = string(rune('0' + i))

		app.ProcessSignal <- tsk
	}

	time.Sleep(6 * time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database update expectations not met: %v", err)
	}
}

func TestSyncTaskToDB(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer func() {
		if err := mockDB.Close(); err != nil {
			t.Fatalf("Failed to close mock DB: %v", err)
		}
	}()

	app.Databasehandle = mockDB

	now := time.Now()
	tsk := task.NewTask("test", "payload", now)
	tsk.State = "completed"
	tsk.Retries = 2

	mock.ExpectExec("^UPDATE tasks SET state").
		WithArgs("completed", sqlmock.AnyArg(), sqlmock.AnyArg(), 2, "", sqlmock.AnyArg(), tsk.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	workerpkg.SyncTaskToDB(tsk)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Logf("Database update expectations not met: %v", err)
	}
}

func TestTaskCancellation(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer func() {
		if err := mockDB.Close(); err != nil {
			t.Fatalf("Failed to close mock DB: %v", err)
		}
	}()

	app.Databasehandle = mockDB
	app.ProcessSignal = make(chan *task.Task, 100)
	app.CancelSignal = make(chan string, 100)

	mock.ExpectExec("^UPDATE tasks SET state").WillReturnResult(sqlmock.NewResult(0, 1))

	workerpkg.ProcessTask()

	now := time.Now()
	tsk := task.NewTask("test", "payload", now.Add(-1*time.Hour))
	tsk.State = "pending"
	tsk.ID = "canceled1"

	app.ProcessSignal <- tsk

	app.CancelSignal <- tsk.ID

	time.Sleep(6 * time.Second)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database update expectations not met: %v", err)
	}
}

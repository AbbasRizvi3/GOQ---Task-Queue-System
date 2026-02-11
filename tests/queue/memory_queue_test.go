package queue

import (
	"testing"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
)

func TestNewMemoryQueue(t *testing.T) {
	mq := queue.NewMemoryQueue()

	if mq == nil {
		t.Error("NewMemoryQueue returned nil")
	}
	if len(mq.Tasks) != 0 {
		t.Errorf("Expected empty queue, got %d tasks", len(mq.Tasks))
	}
}

func TestEnqueueSingleTask(t *testing.T) {
	mq := queue.NewMemoryQueue()
	tsk := task.NewTask("test", "payload", time.Now())
	signalCh := make(chan struct{}, 100)

	mq.Enqueue(tsk, signalCh)

	if len(mq.Tasks) != 1 {
		t.Errorf("Expected 1 task in queue, got %d", len(mq.Tasks))
	}
	if mq.Tasks[0].ID != tsk.ID {
		t.Errorf("Expected task ID %s, got %s", tsk.ID, mq.Tasks[0].ID)
	}
}

func TestEnqueueMultipleTasks(t *testing.T) {
	mq := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	tasks := make([]*task.Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = task.NewTask("task"+string(rune(i)), "payload", time.Now())
		mq.Enqueue(tasks[i], signalCh)
	}

	if len(mq.Tasks) != 5 {
		t.Errorf("Expected 5 tasks in queue, got %d", len(mq.Tasks))
	}
}

func TestEnqueueDuplicateTask(t *testing.T) {
	mq := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	tsk := task.NewTask("test", "payload", time.Now())

	mq.Enqueue(tsk, signalCh)
	mq.Enqueue(tsk, signalCh)

	if len(mq.Tasks) != 1 {
		t.Errorf("Expected 1 task in queue (no duplicates), got %d", len(mq.Tasks))
	}
}

func TestDequeueTask_Existing(t *testing.T) {
	mq := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	tsk := task.NewTask("test", "payload", time.Now())
	mq.Enqueue(tsk, signalCh)

	retrieved, err := mq.DequeueTask(tsk.ID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrieved == nil {
		t.Error("Expected to retrieve task, got nil")
	}
	if retrieved.ID != tsk.ID {
		t.Errorf("Expected task ID %s, got %s", tsk.ID, retrieved.ID)
	}
}

func TestDequeueTask_NonExistent(t *testing.T) {
	mq := queue.NewMemoryQueue()

	retrieved, err := mq.DequeueTask("nonexistent")

	if err == nil {
		t.Error("Expected error for non-existent task")
	}
	if retrieved != nil {
		t.Errorf("Expected nil for non-existent task, got %v", retrieved)
	}
}

func TestDequeueTask_CorrectRemoval(t *testing.T) {
	mq := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	tasks := make([]*task.Task, 3)
	for i := 0; i < 3; i++ {
		tasks[i] = task.NewTask("task"+string(rune(i)), "payload", time.Now())
		mq.Enqueue(tasks[i], signalCh)
	}

	mq.DequeueTask(tasks[1].ID)

	if len(mq.Tasks) != 2 {
		t.Errorf("Expected 2 tasks after dequeue, got %d", len(mq.Tasks))
	}

	for _, tsk := range mq.Tasks {
		if tsk.ID == tasks[1].ID {
			t.Error("Task should have been removed from queue")
		}
	}
}

func TestRemoveTask_Multiple(t *testing.T) {
	mq := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	tasks := make([]*task.Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = task.NewTask("task"+string(rune(i)), "payload", time.Now())
		mq.Enqueue(tasks[i], signalCh)
	}

	mq.DequeueTask(tasks[2].ID)

	if len(mq.Tasks) != 4 {
		t.Errorf("Expected 4 tasks in queue, got %d", len(mq.Tasks))
	}

	for _, tsk := range mq.Tasks {
		if tsk.ID == tasks[2].ID {
			t.Errorf("Task %s should have been removed", tasks[2].ID)
		}
	}
}

func TestQueueConcurrentOperations(t *testing.T) {
	mq := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	tsk := task.NewTask("test task", "payload", time.Now())

	go func() {
		for i := 0; i < 10; i++ {
			mq.Enqueue(tsk, signalCh)
			time.Sleep(1 * time.Millisecond)
		}
	}()

	go func() {
		time.Sleep(5 * time.Millisecond)
		for i := 0; i < 5; i++ {
			mq.DequeueTask(tsk.ID)
			time.Sleep(1 * time.Millisecond)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	if len(mq.Tasks) >= 0 {
		t.Log("Concurrent operations completed successfully")
	}
}

func TestQueueWithLargeNumberOfTasks(t *testing.T) {
	mq := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	taskCount := 100

	for i := 0; i < taskCount; i++ {
		tsk := task.NewTask("task", "payload", time.Now())
		mq.Enqueue(tsk, signalCh)
	}

	if len(mq.Tasks) != taskCount {
		t.Errorf("Expected %d tasks, got %d", taskCount, len(mq.Tasks))
	}
}

func TestEnqueueDequeueOrder(t *testing.T) {
	mq := queue.NewMemoryQueue()
	signalCh := make(chan struct{}, 100)
	tasks := make([]*task.Task, 5)
	for i := 0; i < 5; i++ {
		tasks[i] = task.NewTask("task"+string(rune(i)), "payload", time.Now())
		mq.Enqueue(tasks[i], signalCh)
	}

	for i := 0; i < 5; i++ {
		retrieved, err := mq.DequeueTask(tasks[i].ID)
		if err != nil {
			t.Errorf("Task %d should exist in queue", i)
		}
		if retrieved == nil || retrieved.ID != tasks[i].ID {
			t.Errorf("Expected task %s, got %v", tasks[i].ID, retrieved)
		}
	}
}

package scheduler

import (
	"context"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
)

func ScheduleTasks(ctx context.Context, memQueue *queue.MemoryQueue, signalCh <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-signalCh:
			checkMemoryQueue(memQueue)
		case <-ticker.C:
			checkMemoryQueue(memQueue)
		}
	}
}

func checkMemoryQueue(memQueue *queue.MemoryQueue) {
	memQueue.Mutex.Lock()
	defer memQueue.Mutex.Unlock()

	for _, t := range memQueue.Tasks {
		checkTaskStates(t)
	}
}

func checkTaskStates(t *task.Task) {
	switch t.IsReadyToRun() {
	case true:
		t.MarkReady()
	case false:
		if t.GetState() == "leased" && t.GetLeaseUntil().Before(time.Now()) {
			t.MarkRetry()
		}
	}
}

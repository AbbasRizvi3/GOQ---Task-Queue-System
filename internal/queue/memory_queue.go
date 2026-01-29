package queue

import (
	"sync"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
)

const (
	taskQueueBufferSize = 100
)

type MemoryQueue struct {
	Tasks []*task.Task
	Mutex sync.Mutex
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		Tasks: make([]*task.Task, 0, taskQueueBufferSize),
	}
}

func appendTask(t *task.Task, mq *MemoryQueue) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	mq.Tasks = append(mq.Tasks, t)
}

func (mq *MemoryQueue) Enqueue(t *task.Task, signalCh chan struct{}) {
	appendTask(t, mq)
	t.Mu.Lock()
	defer t.Mu.Unlock()
	select {
	case signalCh <- struct{}{}:
	default:
	}

	runAt := t.RunAt
	if !t.NextRunAt.IsZero() {
		runAt = t.NextRunAt
	}
	if !runAt.IsZero() && runAt.After(time.Now()) {
		go func(d time.Duration) {
			time.Sleep(d)
			select {
			case signalCh <- struct{}{}:
			default:
			}
		}(time.Until(runAt))
	}
}

func (mq *MemoryQueue) DequeueTask(taskID string) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	for i, t := range mq.Tasks {
		if t.ID == taskID {
			mq.Tasks = append(mq.Tasks[:i], mq.Tasks[i+1:]...)
			return
		}
	}
}

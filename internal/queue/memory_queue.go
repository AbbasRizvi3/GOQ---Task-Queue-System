package queue

import (
	"fmt"
	"sync"

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

func (mq *MemoryQueue) Enqueue(t *task.Task, signalCh chan struct{}) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	for _, existing := range mq.Tasks {
		if existing.ID == t.ID {
			return
		}
	}
	mq.Tasks = append(mq.Tasks, t)
}

func (mq *MemoryQueue) DequeueTask(taskID string) (*task.Task, error) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()

	for i := len(mq.Tasks) - 1; i >= 0; i-- {
		if mq.Tasks[i].ID == taskID {
			t := mq.Tasks[i]

			mq.Tasks = append(mq.Tasks[:i], mq.Tasks[i+1:]...)

			return t, nil
		}
	}

	return nil, fmt.Errorf("task with ID %s not found", taskID)
}

func (mq *MemoryQueue) RemoveTask(taskID string) error {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()

	for i := len(mq.Tasks) - 1; i >= 0; i-- {
		if mq.Tasks[i].ID == taskID {
			mq.Tasks = append(mq.Tasks[:i], mq.Tasks[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("task with ID %s not found", taskID)
}

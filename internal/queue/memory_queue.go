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
	tasks []*task.Task
	Mutex sync.Mutex
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		tasks: make([]*task.Task, taskQueueBufferSize),
	}
}

func (mq *MemoryQueue) Enqueue(t *task.Task) error {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	size := len(mq.tasks)
	mq.tasks = append(mq.tasks, t)
	if len(mq.tasks) > size {
		return nil
	}
	return fmt.Errorf("failed to enqueue task")
}

func (mq *MemoryQueue) Dequeue() (*task.Task, error) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	if len(mq.tasks) == 0 {
		return nil, fmt.Errorf("no tasks in queue")
	}
	t := mq.tasks[0]
	mq.tasks = mq.tasks[1:]
	return t, nil
}

func (mq *MemoryQueue) GetTask() (*task.Task, error) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	if len(mq.tasks) == 0 {
		return nil, fmt.Errorf("no tasks in queue")
	}
	t := mq.tasks[0]
	return t, nil
}

package queue

import (
	"sync"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
)

type ResultQueue struct {
	Tasks []*task.Task
	Mutex sync.Mutex
}

func NewResultQueue() *ResultQueue {
	return &ResultQueue{
		Tasks: make([]*task.Task, 0),
	}
}

func (rq *ResultQueue) Enqueue(t *task.Task) {
	rq.Mutex.Lock()
	defer rq.Mutex.Unlock()
	rq.Tasks = append(rq.Tasks, t)
}

func (rq *ResultQueue) GetAllTasks() []*task.Task {
	rq.Mutex.Lock()
	defer rq.Mutex.Unlock()
	return rq.Tasks
}

func (rq *ResultQueue) GetTaskByID(taskID string) *task.Task {
	rq.Mutex.Lock()
	defer rq.Mutex.Unlock()
	for _, t := range rq.Tasks {
		if t.ID == taskID {
			return t
		}
	}
	return nil
}

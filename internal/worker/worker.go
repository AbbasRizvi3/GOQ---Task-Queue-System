package worker

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
)

const (
	leaseTime = 8 * time.Second
)

func ProcessTask(memoryQueue *queue.MemoryQueue, maxWorkers int, ctx context.Context) {
	for range maxWorkers {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				memoryQueue.Mutex.Lock()
				t, err := findReadyTask(memoryQueue)
				if err != nil {
					memoryQueue.Mutex.Unlock()
					time.Sleep(100 * time.Millisecond)
					continue
				}
				t.MarkLeased(int(leaseTime))
				memoryQueue.Mutex.Unlock()

				fmt.Printf("Worker processing task %s\n", t.ID)

				time.Sleep(3 * time.Second)
				if rand.Intn(2) == 0 {
					t.MarkCompleted()
					fmt.Printf("Task %s completed\n", t.ID)
					memoryQueue.DequeueTask(t.ID)
				} else {
					t.MarkFailed("Fatal: Task failed due to random error")
				}
			}
		}()
	}
}

func findReadyTask(mq *queue.MemoryQueue) (*task.Task, error) {
	for _, t := range mq.Tasks {
		t.Mu.Lock()
		defer t.Mu.Unlock()
		isReady := t.State == "ready"

		if isReady {
			return t, nil
		}
	}
	return nil, fmt.Errorf("fatal: no ready tasks")
}

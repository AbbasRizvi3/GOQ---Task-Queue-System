package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/worker"
)

func ScheduleTasks(taskChannel chan struct{}, workerCount int, incrementActiveWorkers func(), activeWorkers int, decrementActiveWorkers func()) {

	go func() {
		for range app.TaskChannel {
			fmt.Println("Task received in TaskChannel")
			if activeWorkers < workerCount {
				incrementActiveWorkers()
				ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
				defer cancel()
				go func() {
					defer decrementActiveWorkers()
					select {
					case <-ctx.Done():
						fmt.Println("Worker finished task or timed out")
					default:
						worker.ProcessTask()
					}
				}()
			}
		}
	}()

}

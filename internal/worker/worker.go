package worker

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/gin-gonic/gin"
)

const (
	maxWorkers = 5
)

var workerCount = 0
var mu sync.Mutex

func increaseWorkers() bool {
	mu.Lock()
	defer mu.Unlock()
	if workerCount >= maxWorkers {
		return false
	}
	workerCount++
	return true
}

func decreaseWorkers() {
	mu.Lock()
	defer mu.Unlock()
	workerCount--
}

func ProcessTask(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker stopped")
				return
			case taskBuffer := <-app.ProcessSignal:
				if !increaseWorkers() {
					fmt.Println("Max workers reached, cannot process more tasks at the moment")
					continue
				}
				t, err := checkTask(taskBuffer)
				if err != nil {
					fmt.Println("Task is not ready to run, skipping for now:", err)
					decreaseWorkers()
					continue
				}

				fmt.Println("increasing workers")
				fmt.Println("no of workers: ", workerCount)

				go func(task *task.Task) {
					defer func() {
						fmt.Println("no of workers: ", workerCount)
						fmt.Println("decreasing workers")
						decreaseWorkers()
					}()
					assignedLease := task.LeaseUntil
					fmt.Printf("Worker processing task %s\n", task.ID)
					time.Sleep(5 * time.Second)

					var currentState string
					var currentLease time.Time
					err := app.Databasehandle.QueryRow(
						"SELECT state, lease_until FROM tasks WHERE id = $1",
						task.ID,
					).Scan(&currentState, &currentLease)

					if err != nil {
						fmt.Printf("Error checking task status: %v\n", err)
						return
					}

					if currentState == "canceled" || !currentLease.Equal(assignedLease) {
						fmt.Printf("Worker for task %s lost ownership (State: %s). Aborting.\n", task.ID, currentState)
						return
					}

					if rand.Intn(100) < 70 {
						task.MarkCompleted()
						SyncTaskToDB(task)
						fmt.Printf("Task %s completed\n", task.ID)
						return
					} else {
						task.MarkFailed("Fatal: Task failed due to random error")
						SyncTaskToDB(task)
						switch task.State {
						case "dead":
							fmt.Printf("Dead task %s removed from memory queue\n", task.ID)
						case "retry":
							fmt.Printf("Retry task %s removed from memory queue (will be re-fetched when next_run_at is due)\n", task.ID)
						}
						return
					}
				}(t)

			}
		}

	}()
}

func SyncTaskToDB(t *task.Task) {
	var nextRunAtDB interface{} = nil
	if !t.NextRunAt.IsZero() {
		nextRunAtDB = t.NextRunAt
	}
	_, err := app.Databasehandle.Exec("UPDATE tasks SET state = $1, next_run_at = $2, lease_until = $3, retries = $4, error = $5, updated_at = $6 WHERE id = $7",
		t.State, nextRunAtDB, t.LeaseUntil, t.Retries, t.Error, time.Now(), t.ID)
	if err != nil {
		fmt.Println("Error syncing task to DB:", err)
	} else {
		app.WebsocketChannelManager.BroadcastJSON(t.UserID, gin.H{
			"type":    "TASK_UPDATED",
			"payload": t,
		})
		fmt.Printf("Synced task %s to DB with state %s\n", t.ID, t.State)
	}
}

func checkTask(t *task.Task) (*task.Task, error) {
	if t.State == "leased" {
		return t, nil
	}
	return nil, fmt.Errorf("task %s is not in leased state (actual: %s)", t.ID, t.State)
}

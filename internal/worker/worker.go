package worker

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
	"github.com/gin-gonic/gin"
)

const (
	leaseTime  = 8 * time.Second
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

func getWorkerCount() int {
	mu.Lock()
	defer mu.Unlock()
	return workerCount
}

func ProcessTask(memoryQueue *queue.MemoryQueue) {
	go func() {
		for range app.ProcessSignal {
			for {
				if !increaseWorkers() {
					fmt.Println("Max workers reached, cannot process more tasks at the moment")
					break
				}
				t, err := findReadyTask(memoryQueue)
				if err != nil {
					fmt.Println("Error finding ready task:", err)
					decreaseWorkers()
					break
				}

				fmt.Println("increasing workers")
				fmt.Println("no of workers: ", workerCount)

				go func(task *task.Task) {

					defer func() {
						fmt.Println("no of workers: ", getWorkerCount())
						fmt.Println("decreasing workers")
						decreaseWorkers()
					}()
					fmt.Printf("Worker processing task %s\n", task.ID)
					time.Sleep(5 * time.Second)
					select {
					case cancelID := <-app.CancelSignal:
						if cancelID == task.ID {
							fmt.Printf("Received cancel signal for task %s, marking as canceled\n", task.ID)
							task.MarkCanceled()
							SyncTaskToDB(task)
							_, err := memoryQueue.DequeueTask(task.ID)
							if err != nil {
								fmt.Printf("Error removing canceled task %s from memory queue: %v\n", task.ID, err)
							} else {
								fmt.Printf("Canceled task %s removed from memory queue\n", task.ID)
							}
							return
						}
					default:
					}
					if task.GetState() == "canceled" {
						fmt.Printf("Task %s was canceled, skipping processing\n", task.ID)
						_, err := memoryQueue.DequeueTask(task.ID)
						if err != nil {
							fmt.Printf("Error removing canceled task %s from memory queue: %v\n", task.ID, err)
						} else {
							fmt.Printf("Canceled task %s removed from memory queue\n", task.ID)
						}
						return
					}
					if rand.Intn(100) < 70 {
						task.MarkCompleted()
						SyncTaskToDB(task)
						_, err := memoryQueue.DequeueTask(task.ID)
						if err != nil {
							fmt.Printf("Error removing completed task %s from memory queue: %v\n", task.ID, err)
						} else {
							fmt.Printf("Completed task %s removed from memory queue\n", task.ID)
						}
						fmt.Printf("Task %s completed\n", task.ID)
						return
					} else {
						task.MarkFailed("Fatal: Task failed due to random error")
						SyncTaskToDB(task)
						_, err := memoryQueue.DequeueTask(task.ID)
						if err != nil {
							fmt.Printf("Error removing failed task %s from memory queue: %v\n", task.ID, err)
						} else {
							switch task.State {
							case "dead":
								fmt.Printf("Dead task %s removed from memory queue\n", task.ID)
							case "retry":
								fmt.Printf("Retry task %s removed from memory queue (will be re-fetched when next_run_at is due)\n", task.ID)
							}
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

func findReadyTask(mq *queue.MemoryQueue) (*task.Task, error) {
	mq.Mutex.Lock()
	defer mq.Mutex.Unlock()
	for _, t := range mq.Tasks {
		if t.IsReadyToRun() {
			t.MarkLeased(int(leaseTime.Seconds()))
			SyncTaskToDB(t)
			return t, nil
		}
	}
	return nil, fmt.Errorf("fatal: no ready tasks")
}

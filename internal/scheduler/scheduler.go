package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
)

const (
	batchSize = 20
)

func ScheduleTasks(ctx context.Context, memQueue *queue.MemoryQueue) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Println("scheduler started")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Scheduler stopped")
			return

		case <-ticker.C:
			fetchAndCacheTasks(memQueue)
		}
	}
}
func fetchAndCacheTasks(memQueue *queue.MemoryQueue) {
	rows, err := app.Databasehandle.Query(
		`SELECT id, name, payload, state, run_at, next_run_at, lease_until, max_retries, retries, error, created_at, updated_at
         FROM tasks
         WHERE state IN ('pending', 'retry', 'ready')
           AND (run_at <= NOW() OR (next_run_at IS NOT NULL AND next_run_at <= NOW()))
         ORDER BY run_at
         LIMIT $1`, batchSize)
	if err != nil {
		fmt.Printf("Error querying tasks: %v\n", err)
		return
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.ID, &t.Name, &t.Payload, &t.State, &t.RunAt, &t.NextRunAt, &t.LeaseUntil, &t.MaxRetries, &t.Retries, &t.Error, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			fmt.Printf("Error scanning task: %v\n", err)
			continue
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("Rows error: %v\n", err)
	}

	for _, t := range tasks {
		now := time.Now()
		if !t.RunAt.IsZero() && t.RunAt.After(now) {
			continue
		}
		if !t.NextRunAt.IsZero() && t.NextRunAt.After(now) {
			continue
		}
		memQueue.Enqueue(t, nil)
		select {
		case app.ProcessSignal <- struct{}{}:
		default:
		}
		fmt.Printf("Fetched and enqueued task %s with state %s, run_at %v, next_run_at %v\n", t.ID, t.State, t.RunAt, t.NextRunAt)
	}
}

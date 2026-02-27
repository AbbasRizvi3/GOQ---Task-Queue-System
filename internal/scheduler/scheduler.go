package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
)

const (
	batchSize = 20
)

func ScheduleTasks(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	fmt.Println("scheduler started")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Scheduler stopped")
			return

		case <-ticker.C:
			fetchTasks()
		}
	}
}

func runFetchTasksQuery() (rows *sql.Rows, err error) {
	return app.Databasehandle.Query(
		`UPDATE tasks
         SET state = 'leased', 
             lease_until = NOW() + ($2 || ' seconds')::INTERVAL,
             updated_at = NOW()
         WHERE id IN (
             SELECT id
             FROM tasks
             WHERE state IN ('pending', 'retry')
             AND COALESCE(next_run_at, run_at) <= NOW()
             ORDER BY COALESCE(next_run_at, run_at)
             FOR UPDATE SKIP LOCKED
             LIMIT $1
         )
         RETURNING id, user_id, name, payload, state, run_at, next_run_at, lease_until,
                   max_retries, retries, error, created_at, updated_at;`,
		batchSize, 8)
}

func sendProcessSignal(rows *sql.Rows) {
	for rows.Next() {
		var t task.Task
		var nextRunAtNull sql.NullTime
		err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Payload, &t.State, &t.RunAt, &nextRunAtNull, &t.LeaseUntil, &t.MaxRetries, &t.Retries, &t.Error, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			fmt.Printf("Error scanning task: %v\n", err)
			continue
		}
		if nextRunAtNull.Valid {
			t.NextRunAt = nextRunAtNull.Time
		}
		select {
		case app.ProcessSignal <- &t:
		default:
			fmt.Printf("ProcessSignal channel is full, skipping task %s\n", t.ID)
		}

		fmt.Printf("Fetched and scheduled task %s with state %s, run_at %v, next_run_at %v\n", t.ID, t.State, t.RunAt, t.NextRunAt)

	}
}
func fetchTasks() {
	rows, err := runFetchTasksQuery()
	if err != nil {
		fmt.Printf("Error querying tasks: %v\n", err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Error closing rows: %v\n", err)
		}
	}()

	sendProcessSignal(rows)

}

package main

import (
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/scheduler"
)

func main() {
	scheduler.ScheduleTasks(app.TaskChannel, app.WorkerCount, app.IncrementActiveWorkers, app.ActiveWorkers, app.DecrementActiveWorkers)
}

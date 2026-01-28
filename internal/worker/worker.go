package worker

import (
	"fmt"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
)

// type Worker struct {
// 	ID        string
// 	LeaseTime time.Time
// }

// func NewWorker() *Worker {
// 	return &Worker{
// 		ID: fmt.Sprintf("%v", uuid.New().String()),
// 	}
// }

// func (w *Worker) Start() {
// 	fmt.Printf("Worker %s started\n", w.ID)
// }

// func (w *Worker) Stop() {
// 	fmt.Printf("Worker %s stopped\n", w.ID)
// }

func ProcessTask() {
	task, err := app.MemoryQueue.GetTask()
	if err != nil {
		fmt.Println("No task to process:", err)
		return
	}
	fmt.Printf("Task %s is being processed\n", task.ID)
	if task.IsReadyToRun(){
		fmt.Printf("Task %s is ready to run\n", task.ID)
		
	}

}

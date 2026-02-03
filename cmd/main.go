package main

import (
	"context"
	"fmt"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/db"
	routers "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/router"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/scheduler"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	taskChannelBufferSize = 100
)

var TaskChannel = make(chan struct{}, taskChannelBufferSize)

// var wg sync.WaitGroup

var router *gin.Engine

func Load() error {
	return godotenv.Load("../.env")
}
func main() {
	err := Load()
	if err != nil {
		panic("Error loading .env file")
	}

	app.Databasehandle, err = db.SetupDatabase()
	if err != nil {
		fmt.Printf("Database setup error: %v\n", err)
		panic("Error setting up database")
	} else {
		println("Database set up successfully")
	}
	// for i := 0; i < 3; i++ {
	// 	t := task.NewTask("test-task", "payload", time.Now())
	// 	go app.MemoryQueue.Enqueue(t, app.SignalCh)
	// }

	// for i := 0; i < 80; i++ {
	// 	t := task.NewTask("test-task", "payload", time.Now().Add(5*time.Second))
	// 	go app.MemoryQueue.Enqueue(t, app.SignalCh)
	// }

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// wg.Add(1)
	go scheduler.ScheduleTasks(ctx, &app.MemoryQueue, app.SignalCh)
	// ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go worker.ProcessTask(&app.MemoryQueue, 5, &app.ResultQueue)

	// wg.Wait() // for now wg is not reduced anywhere, so main will wait indefinitely, this will be catered with when routers are added
	router = routers.SetUpRoutes()

	router.Run(":8000")
}

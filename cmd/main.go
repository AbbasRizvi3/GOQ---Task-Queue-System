package main

import (
	"context"
	"fmt"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/db"
	routers "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/router"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/scheduler"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/template"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scheduler.ScheduleTasks(ctx, &app.MemoryQueue)
	defer cancel()

	go worker.ProcessTask(&app.MemoryQueue)

	router = routers.SetUpRoutes()
	template.SetupTemplate(router)

	router.Run(":8000")
}

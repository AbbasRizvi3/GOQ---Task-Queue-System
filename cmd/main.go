package main

import (
	"context"
	"fmt"

	_ "net/http/pprof"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/db"
	routers "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/router"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/scheduler"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/template"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/olahol/melody"
)

var router *gin.Engine

func Load() error {
	_ = godotenv.Load()
	return nil
}
func main() {
	Load()
	var err error
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
	app.MelodyInstance.HandleConnect(func(s *melody.Session) {
		if userID, exists := s.Get("userID"); exists {
			uidStr, ok := userID.(string)
			if !ok {
				uidStr = fmt.Sprintf("%v", userID)
			}

			fmt.Printf("Client connected: %s\n", uidStr)
			app.WebsocketChannelManager.AddClientToChannel(uidStr, s)
		} else {
			fmt.Println("Connection attempt without UserID")
			s.Close()
		}
	})

	app.MelodyInstance.HandleDisconnect(func(s *melody.Session) {
		if userID, exists := s.Get("userID"); exists {
			uidStr := userID.(string)
			app.WebsocketChannelManager.RemoveClientFromChannel(uidStr, s)
		}
	})

	template.SetupTemplate(router)

	router.Run(":8000")
}

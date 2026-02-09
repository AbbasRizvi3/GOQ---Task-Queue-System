package app

import (
	"database/sql"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
)

var MemoryQueue queue.MemoryQueue = *queue.NewMemoryQueue()
var ProcessSignal chan struct{} = make(chan struct{}, 100)
var Databasehandle *sql.DB

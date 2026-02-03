package app

import (
	"database/sql"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
)

var MemoryQueue queue.MemoryQueue = *queue.NewMemoryQueue()
var ResultQueue queue.ResultQueue = *queue.NewResultQueue()

var SignalCh chan struct{} = make(chan struct{}, 1)
var Databasehandle *sql.DB

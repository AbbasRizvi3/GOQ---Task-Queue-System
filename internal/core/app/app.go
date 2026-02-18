package app

import (
	"database/sql"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/socket"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
	"github.com/olahol/melody"
)

var MemoryQueue queue.MemoryQueue = *queue.NewMemoryQueue()
var ProcessSignal chan struct{} = make(chan struct{}, 100)
var CancelSignal chan string = make(chan string, 100)
var Databasehandle *sql.DB
var MelodyInstance *melody.Melody = melody.New()
var WebsocketChannelManager *socket.ChannelManager = socket.NewChannelManager()

package app

import (
	"database/sql"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/socket"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/olahol/melody"
)

var ProcessSignal chan *task.Task = make(chan *task.Task, 100)
var CancelSignal chan string = make(chan string, 100)
var Databasehandle *sql.DB
var MelodyInstance *melody.Melody = melody.New()
var WebsocketChannelManager *socket.ChannelManager = socket.NewChannelManager()

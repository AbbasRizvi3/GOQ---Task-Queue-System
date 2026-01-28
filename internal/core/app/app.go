package app

import (
	"sync"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
)

const (
	taskChannelBufferSize = 100
	WorkerCount           = 5
)

var mutex sync.Mutex
var ActiveWorkers = 0

var TaskChannel = make(chan struct{}, taskChannelBufferSize)
var MemoryQueue queue.MemoryQueue = *queue.NewMemoryQueue()

// var WorkerPool []*worker.Worker = make([]*worker.Worker, WorkerCount)

func IncrementActiveWorkers() {
	mutex.Lock()
	defer mutex.Unlock()
	ActiveWorkers++
}

func DecrementActiveWorkers() {
	mutex.Lock()
	defer mutex.Unlock()
	ActiveWorkers--
}

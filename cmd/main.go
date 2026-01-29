package main

import (
	"context"
	"sync"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/scheduler"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/worker"
)

const (
	taskChannelBufferSize = 100
)

var TaskChannel = make(chan struct{}, taskChannelBufferSize)

var MemoryQueue queue.MemoryQueue = *queue.NewMemoryQueue()
var SignalCh chan struct{} = make(chan struct{}, 1)

var wg sync.WaitGroup

func main() {

	for i := 0; i < 3; i++ {
		t := task.NewTask("test-task", []byte("payload"), time.Now())
		go MemoryQueue.Enqueue(t, SignalCh)
	}

	for i := 0; i < 3; i++ {
		t := task.NewTask("test-task", []byte("payload"), time.Now().Add(5*time.Second))
		go MemoryQueue.Enqueue(t, SignalCh)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(1)
	go scheduler.ScheduleTasks(ctx, &MemoryQueue, SignalCh)
	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go worker.ProcessTask(&MemoryQueue, 5, ctx2)
	

	wg.Wait() // for now wg is not reduced anywhere, so main will wait indefinitely, this will be catered with when routers are added
}

package worker

import (
	"fmt"
	"time"
	"worker-pool-server/pkg/job"
	"worker-pool-server/pkg/utils"
)

type Worker struct {
	Id         int
	JobQueue   chan job.Job
	WorkerPool chan chan job.Job
	QuitChan   chan bool
}

func NewWorker(id int, workerPool chan chan job.Job) *Worker {
	return &Worker{
		Id:         id,
		JobQueue:   make(chan job.Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				fmt.Printf("Worker with id %d Started\n", w.Id)
				fibRes := utils.Fibonacci(job.Number)
				time.Sleep(job.Delay)
				fmt.Printf("Worker with id %d finished with result %d\n", w.Id, fibRes)
			case <-w.QuitChan:
				fmt.Printf("Worker with id %d finished\n", w.Id)
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

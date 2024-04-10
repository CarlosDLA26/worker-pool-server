package dispatcher

import (
	"worker-pool-server/pkg/job"
	"worker-pool-server/pkg/worker"
)

type Dispatcher struct {
	WorkerPool chan chan job.Job
	MaxWorkers int
	JobQueue   chan job.Job
}

func NewDispatcher(jobQueue chan job.Job, workers int) *Dispatcher {
	return &Dispatcher{
		JobQueue:   jobQueue,
		MaxWorkers: workers,
		WorkerPool: make(chan chan job.Job, workers),
	}
}

func (d *Dispatcher) Dispatch() {
	for {
		job := <-d.JobQueue
		go func() {
			workerJobQueue := <-d.WorkerPool
			workerJobQueue <- job
		}()
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := worker.NewWorker(i, d.WorkerPool)
		worker.Start()
	}
	go d.Dispatch()
}

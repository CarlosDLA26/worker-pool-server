package dispatcher

import "worker-pool-server/pkg/job"

type Dispatcher struct {
	WorkerPool chan chan job.Job
	MaxWorkers int
	JobQueue   chan job.Job
}

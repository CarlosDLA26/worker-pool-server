package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"worker-pool-server/pkg/dispatcher"
	"worker-pool-server/pkg/job"
)

func RequestHandler(w http.ResponseWriter, r *http.Request, jobQueue chan job.Job) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Invalid delay", http.StatusBadRequest)
		return
	}

	number, err := strconv.Atoi(r.FormValue("value"))
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	job := job.Job{Name: name, Delay: delay, Number: number}

	jobQueue <- job
	w.WriteHeader(http.StatusCreated)
}

func main() {
	const maxWorkers = 4
	const maxQueueSize = 20

	jobQueue := make(chan job.Job, maxQueueSize)
	dispatcher := dispatcher.NewDispatcher(jobQueue, maxWorkers)
	dispatcher.Run()

	// http://localhost:8000/fib
	http.HandleFunc("/fib", func(w http.ResponseWriter, r *http.Request) {
		RequestHandler(w, r, jobQueue)
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}

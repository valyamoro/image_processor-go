package worker

import (
	"fmt"
	"image_processor-go/internal/core"
	"image_processor-go/internal/processor"
)

type Worker struct {
	ID			int
	JobChan 	<-chan *core.Job
	ResultChan  chan<- string
	Stats		*core.WorkerStats
}

func NewWorker(id int, jobChan <-chan *core.Job, resultChan chan<- string) *Worker {
	return &Worker{
		ID:			id,
		JobChan:	jobChan,
		ResultChan:	resultChan,
		Stats:		&core.WorkerStats{
			ID: id,
		},
	}
}

func (w *Worker) Start() {
	for job := range w.JobChan {
		w.processJob(job)
	}

	fmt.Printf("Worker #%d finished. Jobs done: %d\n", w.ID, w.Stats.JobsDone)
}

func (w *Worker) processJob(job *core.Job) {
	w.Stats.IsBusy = true
	w.Stats.CurrentJob = job

	switch job.Priority {
	case core.PriorityHigh:
		w.Stats.HighPriorityJobs++
	case core.PriorityMid:
		w.Stats.MidPriorityJobs++
	case core.PriorityLow:
		w.Stats.LowPriorityJobs++
	}

	duration, err := processor.ConvertImage(job)

	result := ""
	if err != nil {
		result = fmt.Sprintf(
			"Worker #%d ERROR: %s -> %s: %v",
			w.ID,
			job.InputPath,
			job.OutputPath,
			err,
		)
	} else {
		result = fmt.Sprintf(
			"Worker #%d OK: %s -> %s in %v (priority: %d)",
			w.ID,
			job.InputPath,
			job.OutputPath,
			duration,
			job.Priority,
		)
	}

	select {
	case w.ResultChan <- result:
		// OK
	default:
		fmt.Printf("Result channel full, dropping: %s\n", result)
	}

	w.Stats.JobsDone++
	w.Stats.TotalTime += duration
	w.Stats.IsBusy = false
	w.Stats.CurrentJob = nil
}

func (w *Worker) GetStats() *core.WorkerStats {
	return w.Stats
}

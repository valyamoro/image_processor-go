package core

import "time"

type Job struct {
	ID		    int
	Priority    Priority
	InputPath   string
	OutputPath  string
	FileModTime time.Time
	FileSize 	 int64
	TargetFormat Format
	Quality		 int
}

type WorkerStats struct {
	ID				 int
	JobsDone		 int
	TotalTime		 time.Duration
	IsBusy			 bool
	CurrentJob		 *Job
	HighPriorityJobs int
	MidPriorityJobs  int
	LowPriorityJobs  int
}

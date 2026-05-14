package job

import (
	"fmt"
	"sync"
)

type Job struct {
	number  int
	Command string
	Status  string
}

func (j Job) String() string {
	mostRecent := ""
	if j.number == mostRecentJobNumber {
		mostRecent = "+"
	}
	if j.number == mostRecentJobNumber-1 {
		mostRecent = "-"
	}

	running := "&"
	if j.Status != "Running" {
		running = ""
	}

	return fmt.Sprintf("[%d]%s  %-21s%s %s", j.number, mostRecent, j.Status, j.Command, running)
}

var jobs []Job
var mu sync.Mutex
var mostRecentJobNumber int = 0

func GetAll() []Job {
	mu.Lock()
	defer mu.Unlock()

	ret := make([]Job, len(jobs))
	copy(ret, jobs)
	return ret
}

func MarkDone(jobNumber int) {
	mu.Lock()
	jobs[jobNumber-1].Status = "Done"
	mu.Unlock()
}

func Reap() {
	mu.Lock()
	defer mu.Unlock()

	filtered := jobs[:0]
	for _, job := range jobs {
		if job.Status != "Done" {
			filtered = append(filtered, job)
		}
	}

	jobs = filtered
}

func Add(job Job) (jobNumber int) {
	mostRecentJobNumber += 1
	job.number = mostRecentJobNumber

	mu.Lock()
	jobs = append(jobs, job)
	mu.Unlock()

	return mostRecentJobNumber
}

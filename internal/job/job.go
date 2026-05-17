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
	mu.Lock()
	defer mu.Unlock()

	mostRecent := ""
	if j.number == jobs[len(jobs)-1].number {
		mostRecent = "+"
	}

	if len(jobs) > 1 && j.number == jobs[len(jobs)-2].number {
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
	defer mu.Unlock()

	for i, job := range jobs {
		if job.number == jobNumber {
			jobs[i].Status = "Done"
			break
		}
	}
}

func Reap() (reaped []string) {
	mu.Lock()
	jbs := jobs
	mu.Unlock()

	filtered := jobs[:0]
	for _, job := range jbs {
		if job.Status != "Done" {
			filtered = append(filtered, job)
		} else {
			reaped = append(reaped, job.String())
		}
	}

	jobs = filtered
	return reaped
}

func Add(job Job) (jobNumber int) {
	mostRecentJobNumber += 1
	job.number = mostRecentJobNumber

	mu.Lock()
	jobs = append(jobs, job)
	mu.Unlock()

	return mostRecentJobNumber
}

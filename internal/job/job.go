package job

import (
	"fmt"
	"maps"
	"slices"
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

	last := ""
	switch j.number {
	case mostRecentNumber:
		last = "+"
	case secondMostRecentNumber:
		last = "-"
	}

	running := "&"
	if j.Status != "Running" {
		running = ""
	}

	return fmt.Sprintf("[%d]%s  %-21s%s %s", j.number, last, j.Status, j.Command, running)
}

var jobs map[int]Job = make(map[int]Job)
var mu sync.Mutex
var lastJobNumber int = 0
var mostRecentNumber int = 0
var secondMostRecentNumber int = 0

func GetAll() []Job {
	mu.Lock()
	defer mu.Unlock()

	result := make([]Job, 0, len(jobs))
	for _, job := range jobs {
		result = append(result, job)
	}

	slices.SortFunc(result, func(a, b Job) int {
		return a.number - b.number
	})

	return result
}

func MarkDone(jobNumber int) {
	mu.Lock()
	defer mu.Unlock()

	if job, ok := jobs[jobNumber]; ok {
		job.Status = "Done"
		jobs[jobNumber] = job
	}
}

func Reap() (reaped []string) {
	mu.Lock()
	j := maps.Clone(jobs)
	mu.Unlock()

	for number, job := range j {
		if job.Status == "Done" {
			reaped = append(reaped, job.String())

			mu.Lock()
			delete(jobs, number)
			updateRecents()
			mu.Unlock()
		}
	}

	return reaped
}

func updateRecents() {
	mostRecentNumber = 0
	secondMostRecentNumber = 0
	for number := range jobs {
		if number > mostRecentNumber {
			secondMostRecentNumber = mostRecentNumber
			mostRecentNumber = number
		} else if number > secondMostRecentNumber {
			secondMostRecentNumber = number
		}
	}
}

func Add(job Job) (jobNumber int) {
	mu.Lock()
	defer mu.Unlock()

	jobNumber = -1
	for i := range lastJobNumber + 1 {
		if _, ok := jobs[i]; i > 0 && !ok {
			jobNumber = i
		}
	}

	if jobNumber == -1 {
		lastJobNumber += 1
		jobNumber = lastJobNumber
	}

	job.number = jobNumber
	jobs[jobNumber] = job

	secondMostRecentNumber = mostRecentNumber
	mostRecentNumber = jobNumber

	return jobNumber
}

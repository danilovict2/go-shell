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
	return fmt.Sprintf("[%d]+  %-24s%s", j.number, j.Status, j.Command)
}

var jobs []Job
var mu sync.Mutex
var nextJobNumber int = 0

func GetAll() []Job {
	mu.Lock()
	defer mu.Unlock()

	ret := make([]Job, len(jobs))
	copy(ret, jobs)
	return ret
}

func Add(job Job) {
	nextJobNumber += 1
	job.number = nextJobNumber

	mu.Lock()
	jobs = append(jobs, job)
	mu.Unlock()
}

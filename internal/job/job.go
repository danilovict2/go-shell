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
	if j.number == MostRecentJobNumber {
		mostRecent = "+"
	}
	if j.number == MostRecentJobNumber - 1 {
		mostRecent = "-"
	}

	return fmt.Sprintf("[%d]%s  %-24s%s", j.number, mostRecent, j.Status, j.Command)
}

var jobs []Job
var mu sync.Mutex
var MostRecentJobNumber int = 0

func GetAll() []Job {
	mu.Lock()
	defer mu.Unlock()

	ret := make([]Job, len(jobs))
	copy(ret, jobs)
	return ret
}

func Add(job Job) {
	MostRecentJobNumber += 1
	job.number = MostRecentJobNumber

	mu.Lock()
	jobs = append(jobs, job)
	mu.Unlock()
}

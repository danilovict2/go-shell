package completions

import "sync"

var completions map[string][]string = make(map[string][]string)
var mu sync.Mutex

func Add(command, pathToCompletion string) {
	mu.Lock()
	defer mu.Unlock()
	completions[command] = append(completions[command], pathToCompletion)
}

func Remove(command string) {
	mu.Lock()
	defer mu.Unlock()
	delete(completions, command)
}

func Get(command string) []string {
	mu.Lock()
	compl := completions[command]
	mu.Unlock()
	return compl
}

package parameter

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

var variables map[string]string = make(map[string]string)
var mu sync.Mutex

func Get(variable string) (string, error) {
	mu.Lock()
	v, ok := variables[variable]
	mu.Unlock()

	if !ok {
		return "", ErrNotFound
	}

	return v, nil
}

func Set(variable, value string) {
	mu.Lock()
	variables[variable] = value
	mu.Unlock()
}

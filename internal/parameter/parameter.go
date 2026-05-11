package parameter

import (
	"errors"
	"sync"
	"unicode"
)

var ErrNotFound = errors.New("not found")
var ErrInvalidIdentifier = errors.New("not a valid identifier")

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

func Set(variable, value string) error {
	if err := validateName(variable); err != nil {
		return err
	}

	mu.Lock()
	variables[variable] = value
	mu.Unlock()

	return nil
}

func validateName(varName string) error {
	if varName[0] >= '0' && varName[0] <= '9' {
		return ErrInvalidIdentifier
	}

	for _, r := range varName {
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && r != '_' {
			return ErrInvalidIdentifier
		}
	}

	return nil
}

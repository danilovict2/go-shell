package parameter

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

func Get() (string, error) {
	return "", ErrNotFound
}

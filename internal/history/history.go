package history

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var Commands []string
var LastAppendIndexes map[string]int = make(map[string]int)

func History(args []string) (string, error) {
	var (
		limit int = len(Commands)
		err   error
	)

	if len(args) > 0 {
		switch args[0] {
		case "-r":
			if len(args) == 1 {
				return "", errors.New("history: missing path to history file")
			}
			return "", LoadFromFile(args[1])
		case "-a":
			if len(args) == 1 {
				return "", errors.New("history: missing path to history file")
			}

			if err := WriteToFile(args[1], LastAppendIndexes[args[1]], os.O_WRONLY|os.O_APPEND); err != nil {
				return "", err
			}

			LastAppendIndexes[args[1]] = len(Commands)
			return "", nil
		case "-w":
			if len(args) == 1 {
				return "", errors.New("history: missing path to history file")
			}
			return "", WriteToFile(args[1], 0, os.O_WRONLY|os.O_CREATE)
		default:
			limit, err = strconv.Atoi(args[0])
			if err != nil {
				return "", err
			}

			if limit < 0 {
				return "", errors.New("n can't be negative")
			}
		}
	}

	ret := ""
	for i, command := range Commands {
		if command != "" && (len(Commands)-i <= limit) {
			ret += fmt.Sprintf(" %d %s\n", i+1, command)
		}
	}

	return strings.TrimRight(ret, "\n"), nil
}

func LoadFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	fHist, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	for _, command := range strings.Split(string(fHist), "\n") {
		if command != "" {
			Commands = append(Commands, command)
		}
	}

	LastAppendIndexes[file] = len(Commands)
	return nil
}

func WriteToFile(file string, start, flag int) error {
	if file == "" {
		return errors.New("file path must not be empty")
	}

	stdout, err := os.OpenFile(file, flag, 0644)
	if err != nil {
		return err
	}

	for _, command := range Commands[start:] {
		fmt.Fprintln(stdout, command)
	}

	return nil
}

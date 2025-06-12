package history

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var Commands []string
var lastAppendIndex = -1

func History(args []string) string {
	var err error
	offset := 0
	limit := len(Commands)
	writeIndex := true

	switch {
	case len(args) == 0:
		break
	case args[0] == "-r":
		if len(args) < 2 {
			return "missing path to history file"
		}

		LoadFromFile(args[1])
		return ""
	case args[0] == "-a":
		writeIndex = false
		if lastAppendIndex >= 0 {
			offset = lastAppendIndex
		}
		lastAppendIndex = len(Commands)
	case args[0] == "-w":
		writeIndex = false
	default:
		limit, err = strconv.Atoi(args[0])
		if err != nil {
			return err.Error()
		}

		if limit < 0 {
			return "n can't be negative"
		}
	}

	ret := ""
	for i, command := range Commands[offset:] {
		if command != "" && (len(Commands)-i <= limit) {
			if writeIndex {
				ret += fmt.Sprintf(" %d %s\n", i+1, command)
			} else {
				ret += fmt.Sprintf("%s\n", command)
			}
		}
	}

	return strings.TrimRight(ret, "\n")
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

	lastAppendIndex = len(Commands)
	return nil
}

func WriteToFile(file string) error {
	if file == "" {
		return fmt.Errorf("file path must not be empty")
	}

	stdout, err := os.OpenFile(os.Getenv("HISTFILE"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	for _, command := range Commands[lastAppendIndex:] {
		fmt.Fprintln(stdout, command)
	}

	return nil
}

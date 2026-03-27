package parser

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/executable"
)

type Suffix struct {
	Suffix   string
	Trailing string
}

func autocomplete(prefix string) (suffixes []Suffix) {
	suffixes = make([]Suffix, 0)
	if len(prefix) == 0 {
		return suffixes
	}

	for _, command := range command.Builtins {
		if strings.HasPrefix(command, prefix) {
			suffix := Suffix{Suffix: command[len(prefix):], Trailing: " "}
			suffixes = append(suffixes, suffix)
		}
	}

	executables := executable.FindExecutables()
	for _, command := range executables {
		command = filepath.Base(command)
		var suffix string

		if len(command) >= len(prefix) {
			suffix = command[len(prefix):]
		}

		if strings.HasPrefix(command, prefix) && !slices.Contains(suffixes, suffix) {
			suffix := Suffix{Suffix: command[len(prefix):], Trailing: " "}
			suffixes = append(suffixes, command[len(prefix):])
		}
	}

	// Only autocomplete files if they're a part of a command
	if strings.Contains(prefix, " ") {
		prefix = prefix[strings.Index(prefix, " ")+1:]
		suffixes = append(suffixes, autocompleteFilename(prefix)...)
	}

	slices.Sort(suffixes)
	return suffixes
}

func autocompleteFilename(filePrefix string) (suffixes []string) {
	suffixes = make([]string, 0)

	dir, file := filepath.Split(filePrefix)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	filePrefix = file

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), filePrefix) || filePrefix == "" {
			suffix := f.Name()[len(filePrefix):]
			if f.IsDir() {
				suffix = suffix + "/"
			}

			suffixes = append(suffixes, suffix)
		}
	}

	return suffixes
}

func allHaveSamePrefix(suffixes []string) bool {
	for _, suffix := range suffixes {
		if !strings.HasPrefix(suffix, suffixes[0]) {
			return false
		}
	}

	return true
}

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
	IsFile   bool
}

func (s Suffix) String() string {
	return s.Suffix + s.Trailing
}

func autocomplete(prefix string) (suffixes []Suffix) {
	suffixes = make([]Suffix, 0)
	if len(prefix) == 0 {
		return suffixes
	}

	for _, command := range command.Builtins {
		if strings.HasPrefix(command, prefix) {
			suffixes = append(suffixes, Suffix{Suffix: command[len(prefix):], Trailing: " "})
		}
	}

	executables := executable.FindExecutables()
	for _, command := range executables {
		command = filepath.Base(command)
		var suffix string

		if len(command) >= len(prefix) {
			suffix = command[len(prefix):]
		}

		if strings.HasPrefix(command, prefix) && !slices.ContainsFunc(suffixes, func(s Suffix) bool { return s.Suffix == suffix }) {
			suffixes = append(suffixes, Suffix{Suffix: command[len(prefix):], Trailing: " "})
		}
	}

	// Only autocomplete files if they're a part of a command
	if strings.Contains(prefix, " ") {
		prefix = prefix[strings.LastIndex(prefix, " ")+1:]
		suffixes = append(suffixes, autocompleteFilename(prefix)...)
	}

	slices.SortFunc(suffixes, func(s1, s2 Suffix) int { return strings.Compare(s1.Suffix, s2.Suffix) })
	return suffixes
}

func autocompleteFilename(filePrefix string) (suffixes []Suffix) {
	suffixes = make([]Suffix, 0)

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
			suffix := Suffix{Suffix: f.Name()[len(filePrefix):], Trailing: " ", IsFile: true}
			if f.IsDir() {
				suffix.Trailing = "/"
			}

			suffixes = append(suffixes, suffix)
		}
	}

	return suffixes
}

func commonPrefix(suffixes []Suffix) string {
	if len(suffixes) == 0 {
		return ""
	}

	prefix := suffixes[0].Suffix
	for _, s := range suffixes[1:] {
		for !strings.HasPrefix(s.Suffix, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}

	return prefix
}

func allHaveSamePrefix(suffixes []Suffix) bool {
	if len(suffixes) == 0 {
		return false
	}

	prefix := commonPrefix(suffixes)
	return len(prefix) > 0
}

package parser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/completions"
	"github.com/codecrafters-io/shell-starter-go/internal/executable"
)

type Completion struct {
	Prefix   string
	Suffix   string
	Trailing string
}

func (c Completion) String() string {
	return c.Suffix + c.Trailing
}

func autocomplete(prefix string) (completions []Completion) {
	if len(prefix) == 0 {
		return completions
	}

	for _, command := range command.Builtins {
		if strings.HasPrefix(command, prefix) {
			completions = append(completions, Completion{Prefix: prefix, Suffix: command[len(prefix):], Trailing: " "})
		}
	}

	executables := executable.FindExecutables()
	for _, command := range executables {
		command = filepath.Base(command)
		var completion string

		if len(command) >= len(prefix) {
			completion = command[len(prefix):]
		}

		if strings.HasPrefix(command, prefix) && !slices.ContainsFunc(completions, func(s Completion) bool { return s.Suffix == completion }) {
			completions = append(completions, Completion{Prefix: prefix, Suffix: command[len(prefix):], Trailing: " "})
		}
	}

	// Autocomplete with files and output of the completer
	if strings.Contains(prefix, " ") {
		completerCompletions := autocompleteCompleter(prefix)
		completions = append(completions, completerCompletions...)

		if len(completerCompletions) == 0 {
			lastSpace := strings.LastIndex(prefix, " ")
			prefix = prefix[lastSpace+1:]
			completions = append(completions, autocompleteFilename(prefix)...)
		}
	}

	slices.SortFunc(completions, func(s1, s2 Completion) int { return strings.Compare(s1.Suffix, s2.Suffix) })
	return completions
}

func autocompleteFilename(filePrefix string) (completions []Completion) {
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
			completion := Completion{Prefix: filePrefix, Suffix: f.Name()[len(filePrefix):], Trailing: " "}
			if f.IsDir() {
				completion.Trailing = "/"
			}

			completions = append(completions, completion)
		}
	}

	return completions
}

func commonPrefix(completions []Completion) string {
	if len(completions) == 0 {
		return ""
	}

	prefix := completions[0].Suffix
	for _, s := range completions[1:] {
		for !strings.HasPrefix(s.Suffix, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}

	return prefix
}

func autocompleteCompleter(expression string) (c []Completion) {
	tokens := strings.Fields(expression)
	command, word, prevWord := tokens[0], "", ""
	if len(tokens) > 1 {
		word = tokens[len(tokens)-1]
	}

	if len(tokens) > 2 {
		prevWord = tokens[1]
	}

	compl := completions.Get(command)
	for _, completionScript := range compl {
		cmd := exec.Command(completionScript, command, word, prevWord)
		output, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		scriptCompletions := strings.Split(string(output), "\n")
		for _, cmp := range scriptCompletions {
			if cmp != "" && strings.HasPrefix(cmp, word) {
				c = append(c, Completion{Prefix: word, Suffix: cmp[len(word):], Trailing: " "})
			}
		}

	}

	return c
}

package command

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/completions"
	"github.com/codecrafters-io/shell-starter-go/internal/executable"
	"github.com/codecrafters-io/shell-starter-go/internal/history"
	"github.com/codecrafters-io/shell-starter-go/internal/parameter"
)

var Builtins []string = []string{"exit", "echo", "type", "pwd", "cd", "history", "jobs", "complete", "declare"}

type Handler func([]string) (string, error)

var BuiltinHandlers map[string]Handler = map[string]Handler{
	"exit":     exit,
	"echo":     echo,
	"type":     commType,
	"pwd":      pwd,
	"cd":       cd,
	"history":  history.History,
	"jobs":     jobs,
	"complete": complete,
	"declare":  declare,
}

func exit(args []string) (string, error) {
	var (
		exitCode int
		err      error
	)

	if len(args) > 0 {
		exitCode, err = strconv.Atoi(args[0])
		if err != nil {
			return "", err
		}
	}

	history.WriteToFile(os.Getenv("HISTFILE"), history.LastAppendIndexes[os.Getenv("HISTFILE")], os.O_WRONLY|os.O_APPEND)
	os.Exit(exitCode)
	return "", nil
}

func echo(args []string) (string, error) {
	return strings.Join(args, " "), nil
}

func commType(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("type: missing operand")
	}

	isBuiltin := slices.Contains(Builtins, args[0])
	if isBuiltin {
		return fmt.Sprintf("%s is a shell builtin", args[0]), nil
	}

	executableFilePath := executable.GetExecutableFilePath(args[0])
	if executableFilePath != "" {
		return fmt.Sprintf("%s is %s", args[0], executableFilePath), nil
	}

	return fmt.Sprintf("%s: not found", args[0]), nil
}

func pwd(args []string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir, nil
}

func cd(args []string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if len(args) == 0 {
		args = append(args, home)
	}

	if args[0] == "~" {
		args[0] = home
	}

	err = os.Chdir(args[0])
	if err != nil {
		return "", fmt.Errorf("cd: %s: No such file or directory", args[0])
	}

	return "", nil
}

func jobs(args []string) (string, error) {
	return "", nil
}

func complete(args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}

	switch {
	case len(args) >= 2 && args[0] == "-p":
		command := args[1]
		compl := completions.Get(command)
		if len(compl) == 0 {
			return fmt.Sprintf("complete: %s: no completion specification", command), nil
		}

		out := ""
		for _, c := range compl {
			out += fmt.Sprintf("complete -C '%s' %s\n", c, command)
		}

		return out[:len(out)-1], nil
	case len(args) >= 3 && args[0] == "-C":
		completions.Add(args[2], args[1])
	case len(args) >= 2 && args[0] == "-r":
		completions.Remove(args[1])
	}

	return "", nil
}

func declare(args []string) (string, error) {
	switch {
	case len(args) >= 2 && args[0] == "-p":
		v, err := parameter.Get(args[1])
		if err != nil {
			return "", fmt.Errorf("declare: %s: %w", args[1], err)
		}

		return fmt.Sprintf("declare -- %s=\"%s\"", args[1], v), nil
	case len(args) > 0 && strings.Contains(args[0], "="):
		parts := strings.Split(args[0], "=")
		varName, value := parts[0], parts[1]
		err := parameter.Set(varName, value)
		if err != nil {
			return "", fmt.Errorf("declare: `%s': %w", args[0], err)
		}
	}

	return "", nil
}

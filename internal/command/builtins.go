package command

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
	"github.com/codecrafters-io/shell-starter-go/internal/history"
)

var Builtins []string = []string{"exit", "echo", "type", "pwd", "cd", "history"}

type Handler func([]string) string

var BuiltinHandlers map[string]Handler = map[string]Handler{
	"exit":    exit,
	"echo":    echo,
	"type":    commType,
	"pwd":     pwd,
	"cd":      cd,
	"history": history.History,
}

func exit(args []string) string {
	var (
		exitCode int
		err      error
	)

	if len(args) > 0 {
		exitCode, err = strconv.Atoi(args[0])
		if err != nil {
			return err.Error()
		}
	}

	history.WriteToFile(os.Getenv("HISTFILE"))
	os.Exit(exitCode)
	return ""
}

func echo(args []string) string {
	return strings.Join(args, " ")
}

func commType(args []string) string {
	if len(args) != 1 {
		return ""
	}

	isBuiltin := slices.Contains(Builtins, args[0])
	if isBuiltin {
		return fmt.Sprintf("%s is a shell builtin", args[0])
	}

	executableFilePath := executable.GetExecutableFilePath(args[0])
	if executableFilePath != "" {
		return fmt.Sprintf("%s is %s", args[0], executableFilePath)
	}

	return fmt.Sprintf("%s: not found", args[0])
}

func pwd(args []string) string {
	dir, err := os.Getwd()
	if err != nil {
		return err.Error()
	}

	return dir
}

func cd(args []string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return err.Error()
	}

	if len(args) == 0 {
		args = append(args, home)
	}

	if args[0] == "~" {
		args[0] = home
	}

	err = os.Chdir(args[0])
	if err != nil {
		return fmt.Sprintf("cd: %s: No such file or directory", args[0])
	}

	return ""
}

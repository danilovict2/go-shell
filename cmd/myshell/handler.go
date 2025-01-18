package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/executable"
)

type Handler func([]string) string

var BuiltinHandlers map[string]Handler = map[string]Handler{
	"exit": exit,
	"echo": echo,
	"type": commType,
	"pwd":  pwd,
	"cd":   cd,
}

func exit(args []string) string {
	var (
		exitCode int
		err error
	)

	if len(args) > 0 {
		exitCode, err = strconv.Atoi(args[0])
		if err != nil {
			return err.Error()
		}
	}

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
	
	isBuiltin := slices.Contains(command.Builtins, args[0])
	if isBuiltin {
		return fmt.Sprintf("%s is a shell builtin", args[0])
	}

	executableFilePaths := executable.FindExecutableFilePaths(args[0])
	if len(executableFilePaths) > 0 {
		return fmt.Sprintf("%s is %s", args[0], executableFilePaths[0])
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

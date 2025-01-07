package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
)

type Handler func([]string) string

var BuiltinHandlers map[string]Handler = map[string]Handler{
	"exit": exit,
	"echo": echo,
	"type": commType,
}

var Builtins []string = []string{"exit", "echo", "type"}

func exit([]string) string {
	os.Exit(0)
	return ""
}

func echo(args []string) string {
	ret := ""
	for _, arg := range args {
		ret += fmt.Sprintf("%s ", arg)
	}

	return ret
}

func commType(args []string) string {
	isBuiltin := slices.Contains(Builtins, args[0])
	if isBuiltin {
		return fmt.Sprintf("%s is a shell builtin", args[0])
	}

	executableFilePaths := executable.FindExecutableFilePaths(args[0])
	if len(executableFilePaths) > 0 {
		return fmt.Sprintf("%s is %s", args[0], executableFilePaths[0])
	}

	return fmt.Sprintf("%s: not found", args[0])
}

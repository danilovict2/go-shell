package main

import (
	"fmt"
	"os"
)

type Handler func([]string) string

var Types = map[string]string{
	"exit": "builtin",
	"echo": "builtin",
	"type": "builtin",
}

var Handlers = map[string]Handler{
	"exit": exit,
	"echo": echo,
	"type": commType,
}

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
	_type, found := Types[args[0]]
	if !found {
		return fmt.Sprintf("%s: not found", args[0])
	}

	return fmt.Sprintf("%s is a shell %s", args[0], _type)
}

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type Handler func([]string) string

var Handlers map[string]Handler = map[string]Handler{
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

	executableFilePaths := make([]string, 0)
	paths := strings.Split(os.Getenv("PATH"), ":")
	wg := sync.WaitGroup{}

	for _, path := range paths {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filepath.Walk(path, func(fPath string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && info.Name() == args[0] {
					executableFilePaths = append(executableFilePaths, fPath)
				}

				return nil
			})
		}()
	}

	wg.Wait()
	if len(executableFilePaths) > 0 {
		return fmt.Sprintf("%s is %s", args[0], executableFilePaths[0])
	}

	return fmt.Sprintf("%s: not found", args[0])
}

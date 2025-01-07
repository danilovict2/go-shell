package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
	"github.com/codecrafters-io/shell-starter-go/internal/reader"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		reader := reader.New(bufio.NewReader(os.Stdin))

		command, err := reader.Read()
		if err != nil {
			fmt.Fprintln(os.Stdout, "error reading input:", err)
			os.Exit(1)
		}

		handler, isBuiltin := BuiltinHandlers[command.Name]
		if isBuiltin {
			fmt.Fprintln(os.Stdout, handler(command.Args))
			continue
		}

		err = executable.Execute(command)
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s: command not found\n", command.Name)
			continue
		}
	}
}

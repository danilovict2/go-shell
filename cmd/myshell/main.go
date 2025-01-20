package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		parser := parser.New(bufio.NewReader(os.Stdin))

		command, err := parser.ParseInput()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			break
		}

		stdout, stderr, err := command.GetOutputWriters()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		handler, isBuiltin := BuiltinHandlers[command.Name]
		if isBuiltin {
			if output := handler(command.Args); output != "" {
				fmt.Fprintln(stdout, handler(command.Args))
			} else {
				fmt.Print(output)
			}

			continue
		}

		err = executable.Execute(command, stdout, stderr)
		if err != nil {
			if err.Error() == "command not found" {
				fmt.Fprintf(os.Stderr, "%s: command not found\n", command.Name)
			}
		}
	}
}
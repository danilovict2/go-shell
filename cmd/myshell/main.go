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

		commands, err := parser.ParseInput()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			break
		}
		
		stdout, stderr, err := commands[0].GetOutputWriters()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		handler, isBuiltin := BuiltinHandlers[commands[0].Name]
		if isBuiltin {
			if output := handler(commands[0].Args); output != "" {
				fmt.Fprintln(stdout, handler(commands[0].Args))
			} else {
				fmt.Print(output)
			}

			continue
		}

		err = executable.Execute(commands[0], stdout, stderr)
		if err != nil {
			if err.Error() == "command not found" {
				fmt.Fprintf(os.Stderr, "%s: command not found\n", commands[0].Name)
			}
		}
	}
}
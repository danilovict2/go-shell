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
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			break
		}

		stdout, err := command.GetOutputWriter()
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

		err = executable.Execute(command, stdout, os.Stderr)
		if err != nil {
			if err.Error() == "command not found" {
				fmt.Fprintf(os.Stderr, "%s: command not found\n", command.Name)
			}
		}
	}
}
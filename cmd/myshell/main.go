package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to set terminal to raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	parser := parser.New(bufio.NewReader(os.Stdin))

	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := parser.ParseInput()
		if err != nil {
			fmt.Fprint(os.Stderr, "error reading input:", err, "\r\n")
			break
		}

		stdout, stderr, err := command.GetOutputWriters()
		if err != nil {
			fmt.Fprint(os.Stderr, err, "\r\n")
			break
		}

		handler, isBuiltin := BuiltinHandlers[command.Name]
		if isBuiltin {
			if output := handler(command.Args); output != "" {
				fmt.Fprint(stdout, handler(command.Args), "\r\n")
			}
			continue
		}

		output, err := executable.Execute(command)
		if err != nil {
			if err.Error() == "command not found" {
				fmt.Fprintf(stderr, "%s: command not found\r\n", command.Name)
			} else {
				fmt.Fprint(stderr, err, "\r\n")
			}
			continue
		}

		fmt.Fprintf(stdout, "%s\r\n", output)
	}
}
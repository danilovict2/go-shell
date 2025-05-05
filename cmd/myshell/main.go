package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
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

		cmd, err := command.Pipeline(commands)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		stdout, stderr, err := cmd.GetOutputWriters()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		output, err := cmd.GetOutput()
		if err != nil {
			fmt.Fprintln(stderr, err)
		} else if output != "" {
			fmt.Fprintln(stdout, output)
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"

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

		switch command.Name {
		case "exit":
			os.Exit(0)
		case "echo":
			for _, arg := range command.Args {
				fmt.Fprintf(os.Stdout, "%s ", arg)
			}

			fmt.Println()
		default:
			fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
		}
	}
}

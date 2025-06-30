package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/history"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

func main() {
	if os.Getenv("HISTFILE") != "" {
		if err := history.LoadFromFile(os.Getenv("HISTFILE")); err != nil {
			log.Fatal(err)
		}
	}

	defer history.WriteToFile(os.Getenv("HISTFILE"), history.LastAppendIndexes[os.Getenv("HISTFILE")], os.O_WRONLY|os.O_APPEND)

	for {
		fmt.Fprint(os.Stdout, "$ ")
		parser := parser.New(bufio.NewReader(os.Stdin))

		commands, err := parser.ParseInput()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)
			break
		}

		command.Pipeline(commands)
	}

}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stdout, "error reading input:", err)
			os.Exit(1)
		}

		command = strings.Trim(command, "\n")
		switch command {
		case "exit 0":
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stdout, "%s: command not found\n", command)
		}
	}
}

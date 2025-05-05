package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
)

type Command struct {
	Name   string
	Args   []string
	Stdout io.WriteCloser
	Stderr io.WriteCloser
}

func New(name string, args []string) *Command {
	c := &Command{
		Name: name,
		Args: args,
	}

	c.setOutputWriters()
	return c
}

func (c Command) String() string {
	return c.Name
}

func (c *Command) setOutputWriters() {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	var err error

	for i := range len(c.Args) - 1 {
		switch c.Args[i] {
		case "'>'", "'1>'":
			c.Stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'>>'", "'1>>'":
			c.Stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'2>'":
			c.Stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'2>>'":
			c.Stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		}
	}
}

func Pipeline(commands []Command) {
	if len(commands) == 1 {
		commands[0].execute()
	}

	for i := 1; i < len(commands); i++ {

	}
}

func (c *Command) execute() {
	handler, isBuiltin := BuiltinHandlers[c.Name]
	if isBuiltin {
		if output := handler(c.Args); output != "" {
			fmt.Fprintln(c.Stdout, output)
		}
		return
	}

	c.executeNonBuiltin()
}

func (c *Command) executeNonBuiltin() {
	if executable.GetExecutableFilePath(c.Name) == "" {
		fmt.Fprintf(c.Stderr, "%s: command not found\n", c.Name)
		return
	}

	comm := exec.Command(c.Name, c.Args...)
	stdin, err := comm.StdinPipe()
	if err != nil {
		fmt.Fprintln(c.Stderr, err)
		return
	}

	go func() {
		defer stdin.Close()
		//io.WriteString(stdin, c.Input)
	}()

	var stderrBuf bytes.Buffer
	comm.Stderr = &stderrBuf

	out, err := comm.Output()
	if err != nil {
		msg := strings.TrimSpace(stderrBuf.String())
		fmt.Fprintln(c.Stderr, msg)
	}

	if string(out) != "" {
		fmt.Fprintln(c.Stdout, strings.TrimRight(string(out), "\n"))
	}
}
